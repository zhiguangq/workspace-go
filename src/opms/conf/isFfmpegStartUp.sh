#!/bin/bash
ps -ef |grep -v bash |grep  ${1} |grep ${2}
