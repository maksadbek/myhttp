myhttp

* Clone and Compile

```bash
git clone git@github.com:maksadbek/myhttp.git
cd myhttp
go build .
```

* Run

```bash
./myhttp http://www.adjust.com http://google.com
http://www.adjust.com	6a485423e7339014b8ca2da5d2b6dd07
http://google.com	9303883bedea3f5e8215b2d5111b5075

./myhttp adjust.com google.com facebook.com yahoo.com yandex.com twitter.com reddit.com/r/funny reddit.com/r/notfunny baroquemusiclibrary.com
http://yandex.com	16a7562c424db7033a1a9ee1bd39e3a2
http://adjust.com	6a485423e7339014b8ca2da5d2b6dd07
http://google.com	b8872be62cd78ae331390cfc4ed8aff8
http://baroquemusiclibrary.com	8aa8db1e3a5d72fcb90fa83856ebfd13
http://twitter.com	4e3abdad931c08fe25c1c5b789846c74
http://reddit.com/r/notfunny	744fe4d01aa100bf43e812caeac127d3
http://reddit.com/r/funny	9ade3f74bf9bf46c3e9499ce6c19346e
http://facebook.com	1265b39d432126610cb0447f7c152084
http://yahoo.com	1f52a5884621a019b038229d599d064c
```
