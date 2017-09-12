# GiveMeTime
GiveMeTime watch source code.  
GiveMeTime is a wrist-watch that displays realtime information about sun, moon and others time related data. The watch itself is a simple microcontroller that handle a 320x320px display, with months of autonomy. It relies on a bespoke ephemerid to display current informations for day, time, sun, moon, etc.  
Ephemeris data for the GiveMeTime watch are generated by the « watch-synchronizer » device/software.

##### « watch-synchronizer » directory contains source code for :
* GPS acquisition/synchronisation.
* Astronomical computations.
* Ephemeris generation. Ephemeris is written to a file descriptor. It will be use by the watch.

This code is intended to be run on a GPS capable device with (small) computational power. It will generate an __ephemeris__ of few kilo bytes of datas. Ephemeris has to be transmitted to the watch by any mean, ie. shared memory, shared fs, usb, bluetooth, beam, whatever.  
watch-synchronizer is written in GO, thus it should be easily portable.

##### « watch » directory contains source code for :
* Ephemeris acquisition/storage/reading
* Clock logic, ie. oscillator reading, registers updating, etc.
* Display controller
* Power management
* Oscillator adjustment (by mean of radio beam or other technics)

This code is intended to run on a very low power specialized device with very small computational capabilities.  
Proof-of-concept is written in GO and target a virtual SVG display.  

As far as today, we do not know if the _watch-synchronizer_ could fit within a wearable device without draining to much power. For the proof-of-concept, we will run this code on a small dedicated device, or a computer, that will send ephemeris data to the watch by any mean (usb, bluetooth, etc.).  
The watch can run without a new time synchronization for up to an year. It will hold a good accuracy for sun & moon as long as it does not move too much away from its synchronization position ("too much away" means below a hundred kilometers). Instead, a new synchronisation has to be ignited for the watch to get an updated ephemeris.