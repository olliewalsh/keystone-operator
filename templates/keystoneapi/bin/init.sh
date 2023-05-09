#!/bin//bash
#
# Copyright 2020 Red Hat Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may
# not use this file except in compliance with the License. You may obtain
# a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
# WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
# License for the specific language governing permissions and limitations
# under the License.
set -ex

# Secrets are obtained from ENV variables.
export PASSWORD=${AdminPassword:?"Please specify a AdminPassword variable."}
export DBHOST=${DatabaseHost:?"Please specify a DatabaseHost variable."}
export DBUSER=${DatabaseUser:?"Please specify a DatabaseUser variable."}
export DBPASSWORD=${DatabasePassword:?"Please specify a DatabasePassword variable."}
export DB=${DatabaseName:-"keystone"}

SVC_CFG=/etc/keystone/keystone.conf
SVC_CFG_MERGED=/var/lib/config-data/merged/keystone.conf

# expect that the common.sh is in the same dir as the calling script
SCRIPTPATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
. ${SCRIPTPATH}/common.sh --source-only

# Copy default service config from container image as base
cp -a ${SVC_CFG} ${SVC_CFG_MERGED}

# Merge all templates from config CM
for dir in /var/lib/config-data/default;do
    merge_config_dir ${dir}
done

# set secrets
if [[ -e /var/lib/tlsdb-data/ca.crt ]]; then
    DB_QS='?ssl=1&ssl_ca=/var/lib/tlsdb-data/ca.crt'
    if [[ -e /var/lib/tlsdb-data/tls.key && -e /var/lib/tlsdb-data/tls.crt ]]; then
        DBUSER=${DBUSER}_tls
        DB_QS=${DB_QS}'&ssl_key=/var/lib/tlsdb-data/tls.key&ssl_cert=/var/lib/tlsdb-data/tls.crt'
    fi
    DB_URI=mysql+pymysql://${DBUSER}@${DBHOST}/${DB}${DB_QS}
else
    DB_URI=mysql+pymysql://${DBUSER}:${DBPASSWORD}@${DBHOST}/${DB}
fi

crudini --set ${SVC_CFG_MERGED} database connection ${DB_URI}
