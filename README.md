# anderson

*checks your go dependencies for contraband licenses*

![Judge Anderson](http://www.scifibloggers.com/wp-content/uploads/dredd-2012.jpg)

## usage

```
$ anderson
Hold still citizen, scanning dependencies for contraband...

github.com/xoebus/apache                                    CHECKS OUT
github.com/xoebus/copyright                                 CONTRABAND
github.com/xoebus/mit                                       CHECKS OUT

We found questionable material. Citizen, what do you have to say for yourself?

github.com/xoebus/no-license                                NO LICENSE
github.com/xoebus/greylist                                  BORDERLINE
```

## configuration

You can configure *anderson* to be more or less lenient when checking you dependencies. A file called .anderson.yml in the root of your Go package will be checked when you run it.

``` yml
whitelist:
- MIT

greylist:
- Apache

blacklist:
- GPL

exceptions:
- github.com/xoebus/greylist
```
