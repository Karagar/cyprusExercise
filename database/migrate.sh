for i in {1..30};
do
    /opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P $MSSQL_SA_PASSWORD -d master -i migration.sql
    if [ $? -eq 0 ]
    then
        echo "migration.sql completed"
        break
    else
        sleep 1
    fi
done