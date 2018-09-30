### Gorilla: Results Data Keeper

  Simplistic App for storing and quering test reports. Using Postgresql (the only option so far) as RDBMS and D3.js library for WebUI.

  This is meant to be practical and not monstrous.

#### Some notes
    curl -vX POST http://localhost:3000/api/upload \
         -d @rdb-aaa.xml \
         --header "Content-Type: application/xml"
