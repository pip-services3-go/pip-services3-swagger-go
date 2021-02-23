package example_logic

import (
	ccomand "github.com/pip-services3-go/pip-services3-commons-go/commands"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	data "github.com/pip-services3-go/pip-services3-swagger-go/example/data"
)

type DummyController struct {
	commandSet *DummyCommandSet
	entities   []data.Dummy
}

func NewDummyController() *DummyController {
	dc := DummyController{}
	dc.entities = make([]data.Dummy, 0)
	return &dc
}

func (c *DummyController) GetCommandSet() *ccomand.CommandSet {
	if c.commandSet == nil {
		c.commandSet = NewDummyCommandSet(c)
	}
	return &c.commandSet.CommandSet
}

func (c *DummyController) GetPageByFilter(correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (items *data.DummyDataPage, err error) {

	if filter == nil {
		filter = cdata.NewEmptyFilterParams()
	}
	var key string = filter.GetAsString("key")

	if paging == nil {
		paging = cdata.NewEmptyPagingParams()
	}
	var skip int64 = paging.GetSkip(0)
	var take int64 = paging.GetTake(100)

	var result []data.Dummy
	for i := 0; i < len(c.entities); i++ {
		var entity data.Dummy = c.entities[i]
		if key != "" && key != entity.Key {
			continue
		}

		skip--
		if skip >= 0 {
			continue
		}

		take--
		if take < 0 {
			break
		}

		result = append(result, entity)
	}
	var total int64 = (int64)(len(result))
	return data.NewDummyDataPage(&total, result), nil
}

func (c *DummyController) GetOneById(correlationId string, id string) (result *data.Dummy, err error) {
	for i := 0; i < len(c.entities); i++ {
		var entity data.Dummy = c.entities[i]
		if id == entity.Id {
			return &entity, nil
		}
	}
	return nil, nil
}

func (c *DummyController) Create(correlationId string, entity data.Dummy) (result *data.Dummy, err error) {
	if entity.Id == "" {
		entity.Id = cdata.IdGenerator.NextLong()
		c.entities = append(c.entities, entity)
	}
	return &entity, nil
}

func (c *DummyController) Update(correlationId string, newEntity data.Dummy) (result *data.Dummy, err error) {
	for index := 0; index < len(c.entities); index++ {
		var entity data.Dummy = c.entities[index]
		if entity.Id == newEntity.Id {
			c.entities[index] = newEntity
			return &newEntity, nil

		}
	}
	return nil, nil
}

func (c *DummyController) DeleteById(correlationId string, id string) (result *data.Dummy, err error) {
	var entity data.Dummy

	for i := 0; i < len(c.entities); {
		entity = c.entities[i]
		if entity.Id == id {
			if i == len(c.entities)-1 {
				c.entities = c.entities[:i]
			} else {
				c.entities = append(c.entities[:i], c.entities[i+1:]...)
			}
		} else {
			i++
		}
		return &entity, nil
	}
	return nil, nil
}
