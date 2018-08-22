package gateway_controller

import "fmt"

func (goc *gwOperationCtx) validate() error {
	if goc.gw.Spec.Image == "" {
		return fmt.Errorf("gateway image is not specified")
	}
	if goc.gw.Spec.Type == "" {
		return fmt.Errorf("gateway type is not specified")
	}
	if len(goc.gw.Spec.Sensors) <= 0 {
		return fmt.Errorf("no associated sensor with gateway")
	}
	if goc.gw.Spec.ServiceAccountName == "" {
		return fmt.Errorf("no service account specified")
	}
	return nil
}