# Http Queue
Simple http-based queue server.

## Example usage

```python
import requests

requests.put("http://server/queueName", data=dataToPutOnTheQueue)
recievedData = requests.get("http://server/queueName").content

```