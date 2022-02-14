package sap_api_caller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	sap_api_output_formatter "sap-api-integrations-planned-order-reads-rmq-kube/SAP_API_Output_Formatter"
	"strings"
	"sync"

	"github.com/latonaio/golang-logging-library-for-sap/logger"
	"golang.org/x/xerrors"
)

type RMQOutputter interface {
	Send(sendQueue string, payload map[string]interface{}) error
}

type SAPAPICaller struct {
	baseURL      string
	apiKey       string
	outputQueues []string
	outputter    RMQOutputter
	log          *logger.Logger
}

func NewSAPAPICaller(baseUrl string, outputQueueTo []string, outputter RMQOutputter, l *logger.Logger) *SAPAPICaller {
	return &SAPAPICaller{
		baseURL:      baseUrl,
		apiKey:       GetApiKey(),
		outputQueues: outputQueueTo,
		outputter:    outputter,
		log:          l,
	}
}

func (c *SAPAPICaller) AsyncGetPlannedOrder(plannedOrder, material, mRPPlant, plant string, accepter []string) {
	wg := &sync.WaitGroup{}
	wg.Add(len(accepter))
	for _, fn := range accepter {
		switch fn {
		case "Header":
			func() {
				c.Header(plannedOrder)
				wg.Done()
			}()
		case "HeaderMaterialPlant":
			func() {
				c.HeaderMaterialPlant(material, mRPPlant)
				wg.Done()
			}()
		case "ComponentMaterialPlant":
			func() {
				c.ComponentMaterialPlant(material, plant)
				wg.Done()
			}()
		default:
			wg.Done()
		}
	}

	wg.Wait()
}

func (c *SAPAPICaller) Header(plannedOrder string) {
	data, err := c.callPlannedOrderSrvAPIRequirementHeader("PlannedOrderHeader", plannedOrder)
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": data, "function": "PlannedOrderHeader"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(data)
}

func (c *SAPAPICaller) callPlannedOrderSrvAPIRequirementHeader(api, plannedOrder string) ([]sap_api_output_formatter.Header, error) {
	url := strings.Join([]string{c.baseURL, "api_plannedorder/srvd_a2x/sap/plannedorder/0001", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithHeader(req, plannedOrder)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, xerrors.Errorf("API status code %d. API request failed", resp.StatusCode)
	}

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToHeader(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) HeaderMaterialPlant(material, mRPPlant string) {
	data, err := c.callPlannedOrderSrvAPIRequirementHeaderMaterialPlant("PlannedOrderHeader", material, mRPPlant)
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": data, "function": "PlannedOrderHeaderMaterialPlant"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(data)
}

func (c *SAPAPICaller) callPlannedOrderSrvAPIRequirementHeaderMaterialPlant(api, material, mRPPlant string) ([]sap_api_output_formatter.Header, error) {
	url := strings.Join([]string{c.baseURL, "api_plannedorder/srvd_a2x/sap/plannedorder/0001", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithHeaderMaterialPlant(req, material, mRPPlant)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, xerrors.Errorf("API status code %d. API request failed", resp.StatusCode)
	}

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToHeader(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) ComponentMaterialPlant(material, plant string) {
	data, err := c.callPlannedOrderSrvAPIRequirementComponentMaterialPlant("PlannedOrderComponent", material, plant)
	if err != nil {
		c.log.Error(err)
		return
	}
	err = c.outputter.Send(c.outputQueues[0], map[string]interface{}{"message": data, "function": "PlannedOrderComponentMaterialPlant"})
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(data)
}

func (c *SAPAPICaller) callPlannedOrderSrvAPIRequirementComponentMaterialPlant(api, material, plant string) ([]sap_api_output_formatter.Component, error) {
	url := strings.Join([]string{c.baseURL, "api_plannedorder/srvd_a2x/sap/plannedorder/0001", api}, "/")
	req, _ := http.NewRequest("GET", url, nil)

	c.setHeaderAPIKeyAccept(req)
	c.getQueryWithComponentMaterialPlant(req, material, plant)

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, xerrors.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, xerrors.Errorf("API status code %d. API request failed", resp.StatusCode)
	}

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToComponent(byteArray, c.log)
	if err != nil {
		return nil, xerrors.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) setHeaderAPIKeyAccept(req *http.Request) {
	req.Header.Set("APIKey", c.apiKey)
	req.Header.Set("Accept", "application/json")
}

func (c *SAPAPICaller) getQueryWithHeader(req *http.Request, plannedOrder string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("PlannedOrder eq '%s'", plannedOrder))
	req.URL.RawQuery = params.Encode()
}

func (c *SAPAPICaller) getQueryWithHeaderMaterialPlant(req *http.Request, material, mRPPlant string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("Material eq '%s' and MRPPlant eq '%s'", material, mRPPlant))
	req.URL.RawQuery = params.Encode()
}

func (c *SAPAPICaller) getQueryWithComponentMaterialPlant(req *http.Request, material, plant string) {
	params := req.URL.Query()
	params.Add("$filter", fmt.Sprintf("Material eq '%s' and Plant eq '%s'", material, plant))
	req.URL.RawQuery = params.Encode()
}
