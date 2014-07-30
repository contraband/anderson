# anderson

*checks your go dependencies for contraband licenses*

![Judge Anderson](http://www.scifibloggers.com/wp-content/uploads/dredd-2012.jpg)

## usage

![Usage](media/usage.png)

## configuration

You can configure *anderson* to be more or less lenient when checking you dependencies. A file called .anderson.yml in the root of your Go package will be checked when you run it.

``` yml
---
whitelist:
- MIT

greylist:
- Apache

blacklist:
- GPL

exceptions:
- github.com/xoebus/greylist
```
