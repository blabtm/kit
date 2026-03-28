package org.blab.v2k.vcas.consumer;

import lombok.Getter;
import lombok.extern.slf4j.Slf4j;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.nio.ByteBuffer;
import java.nio.CharBuffer;
import java.nio.channels.AsynchronousSocketChannel;
import java.nio.channels.CompletionHandler;
import java.nio.charset.StandardCharsets;
import java.sql.DriverManager;
import java.sql.SQLException;
import java.time.Duration;
import java.util.*;

@Slf4j
public class VcasConsumer {
  private final Properties properties;
  private final Integer queueSize;
  private final AsynchronousSocketChannel socket;
  private final ByteBuffer buffer;

  private Map<String, LinkedList<ConsumerRecord>> queue;
  private Boolean isEmpty;

  @Getter
  private final Set<String> topics;

  public VcasConsumer(Properties properties) throws Exception {
    this.properties = properties;
    this.queueSize = (Integer) properties.getOrDefault(ConsumerConfig.QUEUE_SIZE, 64);
    this.socket = AsynchronousSocketChannel.open();
    this.buffer = ByteBuffer.allocate((Integer) properties.getOrDefault(ConsumerConfig.PACKET_MAX, 256));

    this.topics = new TreeSet<>();
    this.queue = new HashMap<>();
    this.isEmpty = true;

    var addr = (String) properties.getOrDefault(ConsumerConfig.BROKER_ADDR, "localhost");
    var port = ((Integer) properties.getOrDefault(ConsumerConfig.BROKER_PORT, 20041));

    this.socket.connect(new InetSocketAddress(addr, port)).get();
    this.socket.read(buffer, null, new InputHandler());
  }

  public void subscribe(Collection<String> tp) throws SQLException, ClassNotFoundException, InterruptedException {
    var t = listTopics(tp);

    log.debug("subscribe: {}", t);

    synchronized (buffer) {
      for (String topic : t) {
        if (topics.contains(topic)) {
          continue;
        }

        socket.write(ByteBuffer.wrap(String
                .format("name:%s|method:subscr\n", topic)
                .getBytes()));

        topics.add(topic);

        Thread.sleep(100);
      }
    }
  }

  private Set<String> listTopics(Collection<String> tp) throws SQLException, ClassNotFoundException {
    var r = new HashSet<String>();

    try (var conn = DriverManager.getConnection(
            properties.getProperty(ConsumerConfig.DB_ADDR),
            properties.getProperty(ConsumerConfig.DB_USER),
            properties.getProperty(ConsumerConfig.DB_PASS)
    )) {
      String query = "SELECT name FROM channels WHERE name ~ ?";

      try (var stmt = conn.prepareStatement(query)) {
        for (var p : tp) {
          stmt.setString(1, p);

          var result = stmt.executeQuery();

          while (result.next()) {
            r.add(result.getString("name"));
          }

          stmt.clearParameters();
        }
      }
    }

    return r;
  }

  public Map<String, LinkedList<ConsumerRecord>> poll(Duration timeout) throws InterruptedException {
    synchronized (buffer) {
      if (timeout.toMillis() == 0) {
        while (isEmpty) {
          buffer.wait();
        }
      } else {
        var st = System.currentTimeMillis();
        var et = st + timeout.toMillis();
        var ct = 0L;

        while (isEmpty && (ct = System.currentTimeMillis()) < et) {
          buffer.wait(et - ct);
        }

        if (isEmpty) {
          return Map.of();
        }
      }

      var ret = queue;

      queue = new HashMap<>();
      isEmpty = true;

      return ret;
    }
  }

  public void close() throws IOException {
    socket.close();
  }

  public class InputHandler implements CompletionHandler<Integer, Void> {
    @Override
    public void completed(Integer n, Void attachment) {
      log.debug("read: {}", n);

      var prv = buffer.position();

      buffer.clear();
      buffer.mark();

      for (int i = 0; i < prv; ++i) {
        if (buffer.get() == '\n') {
          buffer.limit(buffer.position());
          buffer.reset();

          CharBuffer packet = StandardCharsets.UTF_8.decode(buffer);
          ConsumerRecord record = null;

          log.debug("packet: \"{}\", {}", packet.duplicate().toString().replace("\n", "\\n"), packet.length());

          try {
            record = ConsumerRecord.parse(packet
                    .limit(packet.limit() - 1)
                    .toString());
          } catch (Exception e) {
            log.error("malformed record", e);

            buffer.mark();
            buffer.limit(buffer.capacity());

            continue;
          }

          synchronized (buffer) {
            var que = queue.computeIfAbsent(record.getTopic(), k -> new LinkedList<>());

            if (que.size() == queueSize) {
              que.pollFirst();
            }

            que.add(record);

            isEmpty = false;
            buffer.notify();
          }

          buffer.mark();
          buffer.limit(buffer.capacity());
        }
      }

      buffer.limit(buffer.position());
      buffer.reset();
      buffer.compact();

      var left = StandardCharsets.UTF_8
              .decode(ByteBuffer.wrap(Arrays.copyOf(buffer.array(), buffer.position())))
              .toString();

      log.debug("left: \"{}\", {}", left, left.length());

      socket.read(buffer, null, this);
    }

    @Override
    public void failed(Throwable exc, Void attachment) {
      log.error("could not process input", exc);
    }
  }
}
