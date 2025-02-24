#!/usr/bin/env bash

chifra list $1 | cut -f2,3 | tr '\t' ',' >tests/apps.csv
chifra export --accounting --statements 0xccd7fc08532953676ff801791def07d3617ec712 2>/dev/null | cut -f1,2,3,7,12,15,16,17 | tr '\t' ','  >tests/logs.csv
