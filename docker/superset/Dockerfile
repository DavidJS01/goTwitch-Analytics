FROM apache/superset
# Switching to root to install the required packages
USER root
RUN pip install psycopg2-binary

COPY superset/dashboard_gotwitch.zip gotwitch.zip

USER superset
RUN superset superset fab create-admin \
               --username admin \
               --firstname guest \
               --lastname guest \
               --email admin@superset.com \
               --password guest

RUN superset superset db upgrade
RUN superset superset init


# Example: installing the MySQL driver to connect to the metadata database
# if you prefer Postgres, you may want to use `psycopg2-binary` instead

# Example: installing a driver to connect to Redshift
# Find which driver you need based on the analytics database
# you want to connect to here:
# https://superset.apache.org/installation.html#database-dependencies
# Switching back to using the `superset` user
USER root
USER superset
# RUN superset import-dashboards -p ./gotwitch.zip
