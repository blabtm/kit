package org.blab.v2k.nifi.processors.vcas;

import lombok.extern.slf4j.Slf4j;
import org.apache.nifi.annotation.behavior.WritesAttribute;
import org.apache.nifi.annotation.behavior.WritesAttributes;
import org.apache.nifi.annotation.documentation.CapabilityDescription;
import org.apache.nifi.annotation.lifecycle.OnStopped;
import org.apache.nifi.components.PropertyDescriptor;
import org.apache.nifi.processor.*;
import org.apache.nifi.processor.exception.ProcessException;
import org.apache.nifi.processor.util.StandardValidators;
import org.blab.v2k.vcas.consumer.ConsumerConfig;
import org.blab.v2k.vcas.consumer.VcasConsumer;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.util.List;
import java.util.Properties;
import java.util.Set;

@Slf4j
@WritesAttributes({@WritesAttribute(attribute = "filename", description = "Source VCAS topic.")})
@CapabilityDescription("Subscribes to a set of topics and receives messages from the VCAS broker.")
public class ConsumeVCAS extends AbstractProcessor {
  public static final PropertyDescriptor PROP_BROKER_ADDR = new PropertyDescriptor.Builder()
          .name("Broker Address")
          .description("The VCAS broker address.")
          .required(true)
          .addValidator(StandardValidators.NON_EMPTY_VALIDATOR)
          .build();

  public static final PropertyDescriptor PROP_BROKER_PORT = new PropertyDescriptor.Builder()
          .name("Broker Port")
          .description("The VCAS broker port.")
          .defaultValue("20041")
          .addValidator(StandardValidators.NON_EMPTY_VALIDATOR)
          .build();

  public static final PropertyDescriptor PROP_DB_ADDR = new PropertyDescriptor.Builder()
          .name("Database Address")
          .description("The VCAS supporting database address.")
          .required(true)
          .addValidator(StandardValidators.NON_EMPTY_VALIDATOR)
          .build();

  public static final PropertyDescriptor PROP_DB_USER = new PropertyDescriptor.Builder()
          .name("Database User")
          .description("The VCAS supporting database username.")
          .required(true)
          .addValidator(StandardValidators.NON_EMPTY_VALIDATOR)
          .build();

  public static final PropertyDescriptor PROP_DB_PASS = new PropertyDescriptor.Builder()
          .name("Database Password")
          .description("The VCAS supporting database password.")
          .required(true)
          .addValidator(StandardValidators.NON_EMPTY_VALIDATOR)
          .build();

  public static final PropertyDescriptor PROP_TOPIC_NAME = new PropertyDescriptor.Builder()
          .name("Topic Names")
          .description("Comma separated VCAS topic patterns in PostgreSQL style.")
          .required(true)
          .addValidator(StandardValidators.NON_EMPTY_VALIDATOR)
          .build();

  public static final PropertyDescriptor PROP_BUFFER_SIZE = new PropertyDescriptor.Builder()
          .name("Buffer Size")
          .description("Internal buffer size for messages.")
          .required(true)
          .addValidator(StandardValidators.NON_EMPTY_VALIDATOR)
          .build();

  public static final Relationship REL_MESSAGE = new Relationship.Builder()
          .name("Message")
          .description("The VCAS message output in CSV format (comma separated timestamp and value).")
          .build();

  private static final List<PropertyDescriptor> descriptors = List.of(
          PROP_BROKER_ADDR,
          PROP_BROKER_PORT,
          PROP_DB_ADDR,
          PROP_DB_USER,
          PROP_DB_PASS,
          PROP_BUFFER_SIZE,
          PROP_TOPIC_NAME
  );

  private static final Set<Relationship> relationships = Set.of(REL_MESSAGE);

  private VcasConsumer consumer;

  @Override
  public void onTrigger(ProcessContext processContext, ProcessSession processSession) throws ProcessException {
    if (consumer == null) {
      var properties = new Properties();

      properties.put(ConsumerConfig.BROKER_ADDR, processContext.getProperty(PROP_BROKER_ADDR).getValue());
      properties.put(ConsumerConfig.BROKER_PORT, processContext.getProperty(PROP_BROKER_PORT).asInteger());
      properties.put(ConsumerConfig.DB_ADDR, processContext.getProperty(PROP_DB_ADDR).getValue());
      properties.put(ConsumerConfig.DB_USER, processContext.getProperty(PROP_DB_USER).getValue());
      properties.put(ConsumerConfig.DB_PASS, processContext.getProperty(PROP_DB_PASS).getValue());
      properties.put(ConsumerConfig.PACKET_MAX, processContext.getProperty(PROP_BUFFER_SIZE).asInteger());

      try {
        var topic = processContext.getProperty(PROP_TOPIC_NAME).getValue();

        consumer = new VcasConsumer(properties);
        consumer.subscribe(Set.of(topic.split(",")));

        log.info("consumer started, subscribed at {}", consumer.getTopics());
      } catch (Exception e) {
        throw new ProcessException(e);
      }
    }

    try {
      var records = consumer.poll(Duration.ofSeconds(1));

      for (var q : records.entrySet()) {
        for (var r : q.getValue()) {
          var file = processSession.create();

          try (var w = processSession.write(file)) {
            w.write(String.format("%d,%s", r.getTimestamp(), r.getValue())
                    .getBytes(StandardCharsets.UTF_8));
          }

          processSession.putAttribute(file, "filename", q.getKey());
          processSession.transfer(file, REL_MESSAGE);
        }
      }
    } catch (Exception e) {
      throw new ProcessException(e);
    }
  }

  @OnStopped
  public void onStopped() throws IOException {
    log.info("processor stopped");

    if (consumer != null) {
      consumer.close();
      consumer = null;

      log.info("consumer stopped");
    }
  }

  @Override
  protected void init(ProcessorInitializationContext context) {
    super.init(context);
  }

  @Override
  public Set<Relationship> getRelationships() {
    return relationships;
  }

  @Override
  protected List<PropertyDescriptor> getSupportedPropertyDescriptors() {
    return descriptors;
  }
}
