#!/usr/bin/env python

import sys

from google.protobuf.compiler import plugin_pb2 as plugin
from google.protobuf.descriptor_pb2 import FileDescriptorProto

def process_file(f: FileDescriptorProto, res: plugin.CodeGeneratorResponse):
    java_package = f.options.java_package.replace('.', '/')

    for message in f.message_type:
        n = res.file.add()
        n.name = f'{java_package}/{message.name}OuterClass.java'
        n.insertion_point = f'message_implements:{f.package}.{message.name}'
        n.content = "org.blab.v2k.io.Record,"

def process(
        req: plugin.CodeGeneratorRequest, 
        res: plugin.CodeGeneratorResponse
):
    for f in req.proto_file:
        process_file(f, res)

if __name__ == '__main__':
    req = plugin.CodeGeneratorRequest.FromString(sys.stdin.buffer.read())
    res = plugin.CodeGeneratorResponse()

    process(req, res)

    sys.stdout.buffer.write(res.SerializeToString())
