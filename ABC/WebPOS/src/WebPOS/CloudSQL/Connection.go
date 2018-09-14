package CloudSQL

// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - DEL START
//import (
//	"WebPOS/Common"
//	"crypto/tls"
//	"crypto/x509"
//	"database/sql"
//	"errors"
//	"fmt"
//	"io/ioutil"
//
//	"github.com/go-sql-driver/mysql"
//
//	"github.com/goframework/gf/db"
//
//	"github.com/goframework/gf/cfg"
//)
//
//const (
//	SERVER_CONFIG_FILE = "../conf/server.cfg"
//)
//
//const (
//	CFG_KEY_DB_CLOUDSQL_DRIVER        = "DatabaseCloudSQL.Driver"
//	CFG_KEY_DB_CLOUDSQL_HOST          = "DatabaseCloudSQL.Host"
//	CFG_KEY_DB_CLOUDSQL_PORT          = "DatabaseCloudSQL.Port"
//	CFG_KEY_DB_CLOUDSQL_SERVER        = "DatabaseCloudSQL.Server"
//	CFG_KEY_DB_CLOUDSQL_USER          = "DatabaseCloudSQL.User"
//	CFG_KEY_DB_CLOUDSQL_PWD           = "DatabaseCloudSQL.Pwd"
//	CFG_KEY_DB_CLOUDSQL_SCHEMA        = "DatabaseCloudSQL.DatabaseName"
//	CFG_KEY_DB_CLOUDSQL_SSL_CA_FILE   = "DatabaseCloudSQL.SslCaFile"
//	CFG_KEY_DB_CLOUDSQL_SSL_CERT_FILE = "DatabaseCloudSQL.SslCertFile"
//	CFG_KEY_DB_CLOUDSQL_SSL_KEY_FILE  = "DatabaseCloudSQL.SslKeyFile"
//)
//
//var mCfg cfg.Cfg = cfg.Cfg{}
//var mDBFactory *db.SqlDBFactory
//
//func Connect() (*sql.DB, error) {
//
//	// Server.cfgを読み込み
//	mCfg.Load(SERVER_CONFIG_FILE)
//
//	// Server.cfg内に定義されている、Cloud SQL接続に必要な情報を取得
//	cfDbCloudSQLDriver := mCfg.Str(CFG_KEY_DB_CLOUDSQL_DRIVER, "mysql")
//	cfDbCloudSQLHost := mCfg.Str(CFG_KEY_DB_CLOUDSQL_HOST, "127.0.0.1")
//	cfDbCloudSQLPort := mCfg.Int(CFG_KEY_DB_CLOUDSQL_PORT, 3306)
//	cfDbCloudSQLServer := mCfg.Str(CFG_KEY_DB_CLOUDSQL_SERVER, "")
//
//	cfDbCloudSQLUser := mCfg.Str(CFG_KEY_DB_CLOUDSQL_USER, "root")
//	cfDbCloudSQLPwd := mCfg.Str(CFG_KEY_DB_CLOUDSQL_PWD, "")
//	cfDbCloudSQLName := mCfg.Str(CFG_KEY_DB_CLOUDSQL_SCHEMA, "")
//
//	cfDbCloudSQLSslCaFile := mCfg.Str(CFG_KEY_DB_CLOUDSQL_SSL_CA_FILE, "server-ca_dev2.pem")
//	cfDbCloudSQLSslCertFile := mCfg.Str(CFG_KEY_DB_CLOUDSQL_SSL_CERT_FILE, "client-cert_dev2.pem")
//	cfDbCloudSQLSslKeyFile := mCfg.Str(CFG_KEY_DB_CLOUDSQL_SSL_KEY_FILE, "client-key_dev2.pem")
//
//	if cfDbCloudSQLServer == "" {
//		cfDbCloudSQLServer = fmt.Sprintf("%s:%d", cfDbCloudSQLHost, cfDbCloudSQLPort)
//	}
//
//	caCertPool := x509.NewCertPool()
//	caPem, err := ioutil.ReadFile(cfDbCloudSQLSslCaFile)
//	if err != nil {
//		return nil, err
//	}
//
//	if appOk := caCertPool.AppendCertsFromPEM(caPem); !appOk {
//		return nil, errors.New("証明書プールへの" + cfDbCloudSQLSslCaFile + "の追加に失敗しました。")
//	}
//
//	clientCert := make([]tls.Certificate, 0, 1)
//	certs, err := tls.LoadX509KeyPair(cfDbCloudSQLSslCertFile, cfDbCloudSQLSslKeyFile)
//	if err != nil {
//		return nil, err
//	}
//
//	clientCert = append(clientCert, certs)
//	mysql.RegisterTLSConfig("CloudSQL", &tls.Config{
//		RootCAs:            caCertPool,
//		Certificates:       clientCert,
//		InsecureSkipVerify: true,
//	})
//
//	db, err := sql.Open(cfDbCloudSQLDriver,
//		fmt.Sprintf("%[1]s:%[2]s@(%[3]s)/%[4]s?tls=CloudSQL", cfDbCloudSQLUser, cfDbCloudSQLPwd, cfDbCloudSQLServer, cfDbCloudSQLName))
//	if err == nil {
//		Common.LogOutput("Cloud SQLに接続しました。")
//	}
//
//	return db, err
//
//}
// ASO-5843 MariaDBをCLOUDSQLに引っ越して、Mariaを停止する - DEL END