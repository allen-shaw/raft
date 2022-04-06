package conf

import "github.com/AllenShaw19/raft/raft/entity"

type ConfigurationEntry struct {
	id      *entity.LogId
	conf    *Configuration
	oldConf *Configuration
}

func NewConfigurationEntry() *ConfigurationEntry {
	return &ConfigurationEntry{}
}

func (c *ConfigurationEntry) Id() *entity.LogId {
	return c.id
}

func (c *ConfigurationEntry) SetId(id *entity.LogId) {
	c.id = id
}

func (c *ConfigurationEntry) Conf() *Configuration {
	return c.conf
}

func (c *ConfigurationEntry) SetConf(conf *Configuration) {
	c.conf = conf
}

func (c *ConfigurationEntry) OldConf() *Configuration {
	return c.oldConf
}

func (c *ConfigurationEntry) SetOldConf(oldConf *Configuration) {
	c.oldConf = oldConf
}

func (c *ConfigurationEntry) IsEmpty() bool {
	return c.conf.IsEmpty()
}

func (c *ConfigurationEntry) IsValid() bool {

}

func (c *ConfigurationEntry) IsStable() bool {

}
