# chatbot-load-testing

## Introduction
This application is based on [github/alex/fb-chat-emulator](https://github.com/ss-dev/fb-chat-emulator). The original project imitate a facebook app, I used it to make this project, that simulate many users seding messages to a bot. The idea is check the performance of a bot keeping a conversation.

## How it works?
The application can imitate a user. It sends a request to the chat-bot. And for each chatbot reply, we choice the next message.

So you need:
1. write `config.json` like this [config.sample.json](https://github.com/gabicavalcante/chatbot-load-testing/config.sample.json)
2. run binary `$ chatbot-load-testing` for your config
3. change the `ChatBotURL` to your app URL
4. and test it

## Configuration
In a config file you can set:

`ChatBotURL` - base URL to your chat-bot app.

The app spends some time to work our requests so you can set the time between `RequestTimeMin` and `RequestTimeMax` ms.

Users spend some time to respond - from `ResponsePauseMin` to `ResponsePauseMax` ms.

By using `Rules` you can describe some reaction.
For example, chat-bot app sends a question for a user (say during some quiz):

```
POST https://graph.facebook.com/v2.6/me/messages
{"message": {"attachment": {"type": "template", ... "title": "What is AppURL parameter above?"}}, "recipient": {"id": "***************"}}
```

so we can describe it as

```javascript
"Send": {
    "Name": "step#1",
    "Content": "gostaria de saber o preço"
}
```

there `Name` need only for readable statistic.

And if right answer is `It is address of your app server` our `Response` would be like this:

```javascript
"Receive": {
    "Content": "Olá. Seja bem vindo!"
}
```

and it's a full rule:

```javascript
{
    "Send": {
        "Name": "step#1",
        "Content": "gostaria de saber o preço"
    },
    "Receive": {
        "Content": "Olá. Seja bem vindo!"
    }
},
```

You can write empty `Response` then the app won't send any response but will write statistic for the requests.

## Statistics

WIP

## Flags
If you want to ctelhange app port or path to config or say to enable debug mode you should use `flags`.

By the way, Debug mode can help you to find out what exactly do you send to Facebook.

Use `$ chatbot-load-testing --help` for more information.

## Releases

You can build your own binary by using `go tools` or use prepared builds there [https://github.com/gabicavalcante/chatbot-load-testing/releases](https://github.com/gabicavalcante/chatbot-load-testing/releases)
