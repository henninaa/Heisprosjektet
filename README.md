# Heisprosjektet
Heisprosjektet som er aktuell fra 13 i 3
-----------------------------------------

oppskrift på GOPATH:
cd til rett mappe type "go":

export GOPATH=$(pwd)

ferdig med gopath

opskrift på git setup:

cd til rett mappe

git init

git remote add origin "https://www.github.com/henninaa/Heisprosjektet.git"

Må kanskje gjøre dette: 

git config --global user.email "DIN EMAIL"

git congfig --global user.name "DITT USERNAME"

Hvis du har noen filer må du committe:

git commit -m "NAVNET PÅ COMMITEN"

Hvis det står noe i rødt skriv begge eller en av disse:

git add '*'

git stage '*'

Så, skal du pulle kjører du

git pull origin master

Skal du pushe skriver du

git push origin master

er det problemer med å pushe må du pulle først, godta en merge, komme deg ut av en eventuell melding i terminalen med ctr x, så committe og så pushe

Ferdig satt opp!

obs, "master er navnet på branchen vi bruker. Hvis vi lager en ny branch må en bytte ut master med navnet på den nye branchen.
