# Official [shouldiridemybike.com](http://shouldiridemybike.com) source

This is a google appengine app written in Go ([Go AppEngine SDK](https://devel
opers.google.com/appengine/downloads#Google_App_Engine_SDK_for_Go)). It
decides for you if you should ride your bike or better stay at home. It is
deployed at [shouldiridemybike.com](http://shouldiridemybike.com).

## Usage

If you want to fork it or create your own app, these are the steps to get up
and running:

1. Clone: `git clone https://github.com/maxsz/shouldiridemybike.git`
2. Develop: `dev_appserver.py shouldiridemybike`
3. Register your app with google (see [appengine website](https://developers.google.com/appengine))
4. Deploy: `appcfg.py update shouldiridemybike`

## Acknowledgements

- Data powered by [forecast.io](http://forecast.io)
- Uses [Google Maps API](https://developers.google.com/maps)
- Icon based on [Bicycle](http://thenounproject.com/noun/bicycle/#icon-No54)
- Uses [Add to Homescreen](https://github.com/cubiq/add-to-homescreen) Javascript