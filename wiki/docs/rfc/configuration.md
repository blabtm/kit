# RFC BINP 1: Collider Configuration Representation and Management

## Abstract

This memo is a proposal to collider configuration representation and management.

## Introduction

The following proposal defines a mechanism through which a collider can be
managed, configuration data can be retrieved, and new configuration data can be
uploaded and manipulated.

A key aspect of proposed configuration representation is that it allows to
closely mirror the native functionality of the collider.

## Terminology

The following terminology is inherited from the [RFC6241] and modified to match
domain of the current memo:

- datastore: A conceptual place to store and access information. A datastore
  might be implemented, for example, using files, a database, flash memory
  locations, or combinations thereof.

- configuration data: The set of writable data that is required to transform a
  system from its initial default state into its current state.

- configuration datastore: The datastore holding the complete set of
  configuration data that is required to get a collider from its initial default
  state into a desired operational state.

- running configuration datastore: A configuration datastore holding the
  complete configuration currently active on the collider. The running
  configuration datastore always exists.

- candidate configuration datastore: A configuration datastore that can be
  manipulated without impacting the collider's current configuration and that
  can be committed to the running configuration datastore.

- state data: The additional data on a system that is not configuration data
  such as read-only status information and collected statistics.

## Representation

In the current memo collider is viewed as a single complex device, consisting of
several subsystems.

Collider configuration is represented as a structured collection of text files,
indexed by the version control system. Structure and constraints of the
configuration are defined in the domain model [RFC-BINP-2].
