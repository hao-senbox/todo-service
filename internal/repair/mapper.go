package repair

import (
	"log"
	"todo-service/internal/location"
	"todo-service/internal/shop"
	"todo-service/internal/user"
)

func buildRepairResponse(repair *Repair, reportBy interface{}, repairBy interface{}, locationInfor interface{}, imageReport []string, imageRepair []string, shopItems *shop.RepairItemsSummary) *RepairResponse {

	var loc location.LocationInfor

	if locationInfor != nil {
		switch v := locationInfor.(type) {
		case location.LocationInfor:
			loc = v
		case *location.LocationInfor:
			loc = *v
		default:
			log.Printf("buildRepairResponse: Unexpected location type: %T, value: %+v", locationInfor, locationInfor)
		}
	}

	var reporter user.UserInfor
	if reportBy != nil {
		switch v := reportBy.(type) {
		case user.UserInfor:
			reporter = v
		case *user.UserInfor:
			reporter = *v
		default:
			log.Printf("buildRepairResponse: Unexpected reportBy type: %T, value: %+v", reportBy, reportBy)
		}
	}

	var repairer *user.UserInfor
	if repairBy != nil {
		switch v := repairBy.(type) {
		case user.UserInfor:
			repairer = &v
		case *user.UserInfor:
			repairer = v
		default:
			log.Printf("buildRepairResponse: Unexpected repairBy type: %T, value: %+v", repairBy, repairBy)
		}
	}

	var imageRepairPtrs []*string
	for _, img := range imageRepair {
		if img != "" {
			imageRepairPtrs = append(imageRepairPtrs, &img)
		}
	}

	return &RepairResponse{
		ID:             repair.ID,
		OrganizationID: repair.OrganizationID,
		JobNumber:      repair.JobNumber,
		JobName:        repair.JobName,
		QRCode:         repair.QRCode,
		Location:       loc,
		Status:         repair.Status,
		DateReport:     repair.DateReport,
		ReportBy:       reporter,
		UrgentVote:     repair.UrgentVote,
		Comment:        repair.Comment,
		ImageReport:    imageReport,
		DateRepair:     repair.DateRepair,
		RepairBy:       repairer,
		CommentRepair:  repair.CommentRepair,
		ImageRepair:    imageRepairPtrs,
		ShopItems:      shopItems,
		TotalCost:      repair.TotalCost,
		CreatedAt:      repair.CreatedAt,
		UpdatedAt:      repair.UpdatedAt,
	}
}
