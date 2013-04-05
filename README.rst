Samsung Widget Server
=====================

This project helps to ease the deployment of widgets to a samsung TV in a
development environment by deploying widgets directly from the package path,
without having to worry about keep the packages up to date on the web server, or
having to tweak widgets.xml


What it does
------------

``samsing-widget-server`` scans the package path, and creates a widget
descriptor each time the SmartTV makes a request to update packages

Instalation
-----------

    go get github.com/noonien/samsung-widget-server

Usage
-----

   samsung-widget-server [--widgetPath=<path to widgets>]

