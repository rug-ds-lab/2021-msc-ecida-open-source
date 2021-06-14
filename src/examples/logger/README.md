# API request logger

This module starts two webservers: a webserver that logs POST requests to `/input` and a webserver that allows you to view the messages that were sent on this webserver.

For example, with the pipeline webserver (the one with the `POST /input` endpoint) being run on port `8000`, and the spectator webserver being run on port `8001`:

```
POST :8000/input

{
    "temperature": 68
}

> HTTP 200 OK


POST :8000/input

{
    "temperature": 86
}

> HTTP 200 OK


GET :8001 /

> HTTP 200 OK
>
> "[{\"temperature\":68},{\"temperature\":86}]"
```
