using MySql.Data.MySqlClient;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace DioLED
{
    internal class DatabaseHandler
    {
        string elapsed;
        public DatabaseHandler(string elapsed) { 
            this.elapsed = elapsed;
        }

        //function connects to database using a connectionString
        public void ConnectToDB()
        {
            MySqlConnection conn = new MySqlConnection();
            try
            {
                string connetionString = "server=192.168.0.100;uid=root;pwd=TDC_arch2023;database=diobase";
                conn = new MySqlConnection(connetionString);
                conn.Open();
                insertData(conn);
            }
            catch
            {
                //Console.WriteLine(ex.Message);
                //if connection does not work, database doesn't exist
                //create database and table
                string createDbString = "server=127.0.0.1;uid=root;pwd=TDC_arch2023";
                MySqlConnection dbconn = new MySqlConnection(createDbString);
                dbconn.Open();

                MySqlCommand dbcmd = new MySqlCommand();
                dbcmd.Connection = dbconn;
                dbcmd.CommandText = "CREATE DATABASE diobase";
                dbcmd.ExecuteNonQuery();
                dbcmd.CommandText = "CREATE TABLE diobase.dios(id integer PRIMARY KEY AUTO_INCREMENT, duration VARCHAR(20), whattime timestamp)";
                dbcmd.ExecuteNonQuery();

                dbconn.Close();
                //connect to db again 
                ConnectToDB();
            }
            
        }

        //function inserts data into database
        private void insertData(MySqlConnection conn)
        {
            MySqlCommand cmd = new MySqlCommand();
            cmd.Connection = conn;
            cmd.CommandText = "INSERT INTO diobase.dios (duration, whattime) VALUES('" + elapsed + "', CURRENT_TIMESTAMP());";
            cmd.ExecuteNonQuery();
            closeDB(conn);
        }
        
        //function closes database connection
        private void closeDB(MySqlConnection conn)
        {
            conn.Close();
        }

    }
}
