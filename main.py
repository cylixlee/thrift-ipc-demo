import subprocess
import sys
import time
from typing import IO, final, override

from thrift.protocol.TJSONProtocol import TJSONProtocol
from thrift.transport.TTransport import TTransportBase, TTransportException

from src.client.generated.Hello import Client, HelloRequest


@final
class StdioTransport(TTransportBase):
    _closed: bool
    _stdin: IO[bytes]
    _stdout: IO[bytes]

    def __init__(self, process: subprocess.Popen[bytes]) -> None:
        self._closed = False
        if not process.stdin or not process.stdout:
            raise RuntimeError("process STDIN/STDOUT not captured")
        self._stdin = process.stdin
        self._stdout = process.stdout

    @override
    def isOpen(self) -> bool:
        return not self._closed

    @override
    def open(self) -> None:
        self._closed = False

    @override
    def close(self) -> None:
        self._closed = True

    @override
    def read(self, size: int) -> bytes:
        if not self.isOpen():
            raise TTransportException(TTransportException.NOT_OPEN)
        return self._stdout.read(size)

    @override
    def write(self, buf: bytes) -> None:
        if not self.isOpen():
            raise TTransportException(TTransportException.NOT_OPEN)
        self._stdin.write(buf)

    @override
    def flush(self) -> None:
        if not self.isOpen():
            raise TTransportException(TTransportException.NOT_OPEN)
        self._stdin.flush()


def main() -> None:
    proc = subprocess.Popen("./server", stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=sys.stderr)
    time.sleep(1)
    transport = StdioTransport(proc)
    protocol = TJSONProtocol(transport)

    client = Client(protocol, protocol)
    while True:
        line = input(">>> ")
        if line.strip() == "":
            break
        try:
            resp = client.hello(HelloRequest(line))
            print(resp.msg)
        except KeyboardInterrupt:
            break
    proc.terminate()


if __name__ == "__main__":
    main()
