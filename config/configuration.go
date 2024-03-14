package config

// Configurations exported
type Configurations struct {
	Server   ServerConfigurations
	Database DatabaseConfigurations
	Secure   SecureConfigurations
	Cron     CronConfigurations
}

// ServerConfigurations exported
type ServerConfigurations struct {
	Port int
}

type CronConfigurations struct {
	Interval int
	AggregatedInterval int
}

// ServerConfigurations exported
type SecureConfigurations struct {
	Hash 		string
	CertFilePath string
  	KeyFilePath  string
}

// DatabaseConfigurations exported
type DatabaseConfigurations struct {
	DBName     string
	DBUser     string
	DBPassword string
	DBPort     string
	DBHost     string
}
