package repair

import (
	"context"
	"fmt"
	"log"
	"time"
	"todo-service/internal/location"
	"todo-service/internal/shop"
	"todo-service/internal/uploader"
	"todo-service/internal/user"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RepairService interface {
	CreateRepair(ctx context.Context, req CreateRepairRequest, userID string) (*string, error)
	GetRepairs(ctx context.Context) ([]*RepairResponse, error)
	GetRepairByID(ctx context.Context, id string) (*RepairResponse, error)
	UpdateRepair(ctx context.Context, req UpdateRepairRequest, id string, userID string) error
	DeleteRepair(ctx context.Context, id string, userID string) error

	AssignRepair(ctx context.Context, id string, userID string, req AssignRepairRequest) error
	CompleteRepair(ctx context.Context, id string, req CompleteRepairRequest, userID string) error
}

type repairService struct {
	RepairRepo      RepairRepository
	LocationService location.LocationService
	UserService     user.UserService
	UploaderService uploader.ImageService
	ShopService     shop.ShopService
	Policy          *Policy
}

func NewRepairService(
	RepairRepo RepairRepository,
	LocationService location.LocationService,
	UserService user.UserService,
	UploaderService uploader.ImageService,
	ShopService shop.ShopService,
) RepairService {
	return &repairService{
		RepairRepo:      RepairRepo,
		LocationService: LocationService,
		UserService:     UserService,
		UploaderService: UploaderService,
		ShopService:     ShopService,
		Policy:          NewPolicy(),
	}
}

func (s *repairService) CreateRepair(ctx context.Context, req CreateRepairRequest, userID string) (*string, error) {

	if req.OrganizationID == "" {
		return nil, fmt.Errorf("organization_id is required")
	}

	if req.JobName == "" {
		return nil, fmt.Errorf("job_name is required")
	}

	if req.Location == "" {
		return nil, fmt.Errorf("location is required")
	}

	if req.UrgentVote < 0 {
		return nil, fmt.Errorf("urgent_vote must be greater than 0")
	}

	if req.Comment == "" {
		return nil, fmt.Errorf("comment is required")
	}

	if len(req.ImageReport) == 0 {
		return nil, fmt.Errorf("image_report is required")
	}

	jobCount, err := s.RepairRepo.GetJobCount(ctx, req.OrganizationID)
	if err != nil {
		return nil, err
	}

	id := primitive.NewObjectID()
	qrCode := fmt.Sprintf("SENBOX.ORG[REPAIR]:%s", id.Hex())

	repair := &Repair{
		ID:             id,
		OrganizationID: req.OrganizationID,
		JobNumber:      jobCount + 1,
		QRCode:         qrCode,
		JobName:        req.JobName,
		Location:       req.Location,
		AssignedTo:     nil,
		Status:         "pending",
		UrgentVote:     req.UrgentVote,
		Comment:        req.Comment,
		DateReport:     time.Now(),
		ReportBy:       userID,
		ImageReport:    req.ImageReport,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = s.RepairRepo.CreateRepair(ctx, repair)
	if err != nil {
		return nil, err
	}

	repairID := id.Hex()

	return &repairID, nil
}

func (s *repairService) GetRepairs(ctx context.Context) ([]*RepairResponse, error) {
	repairs, err := s.RepairRepo.GetRepairs(ctx)
	if err != nil {
		return nil, err
	}

	var results []*RepairResponse
	for _, repair := range repairs {
		reportBy, err := s.UserService.GetUserInfor(ctx, repair.ReportBy)
		if err != nil {
			return nil, err
		}

		var repairBy interface{}
		if repair.RepairBy != nil {
			repairBy, err = s.UserService.GetUserInfor(ctx, *repair.RepairBy)
			if err != nil {
				return nil, err
			}
		}

		location, err := s.LocationService.GetLocationByID(ctx, repair.Location)
		if err != nil {
			return nil, err
		}

		imageReport := make([]string, len(repair.ImageReport))
		for i, image := range repair.ImageReport {
			imageKey, err := s.UploaderService.GetImageKey(ctx, image)
			if err != nil {
				imageReport[i] = ""
			} else if imageKey == nil {
				imageReport[i] = ""
			} else {
				imageReport[i] = imageKey.Url
			}
		}

		var imageRepair []string
		if repair.ImageRepair != nil {
			imageRepair = make([]string, len(repair.ImageRepair))
			for i, image := range repair.ImageRepair {
				if image != nil {
					imageKey, err := s.UploaderService.GetImageKey(ctx, *image)
					if err != nil {
						imageRepair[i] = ""
					} else if imageKey == nil {
						imageRepair[i] = ""
					} else {
						imageRepair[i] = imageKey.Url
					}
				} else {
					imageRepair[i] = ""
				}
			}
		} else {
			log.Printf("GetRepairByID: repair.ImageRepair is nil")
		}
		// Get shop items for this repair
		shopItems, err := s.ShopService.GetRepairItems(ctx, repair.ID)
		if err != nil {
			log.Printf("GetRepairs: Error getting shop items for repair %s: %v", repair.ID.Hex(), err)
			// Continue without shop items if error
			shopItems = nil
		}

		results = append(results, buildRepairResponse(repair, reportBy, repairBy, location, imageReport, imageRepair, shopItems))
	}

	return results, nil
}

func (s *repairService) GetRepairByID(ctx context.Context, id string) (*RepairResponse, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("GetRepairByID: Error converting ID to ObjectID: %v", err)
		return nil, err
	}

	repair, err := s.RepairRepo.GetRepairByID(ctx, objectID)
	if err != nil {
		log.Printf("GetRepairByID: Error getting repair from repo: %v", err)
		return nil, err
	}

	reportBy, err := s.UserService.GetUserInfor(ctx, repair.ReportBy)
	if err != nil {
		return nil, err
	}

	var repairBy interface{}
	if repair.RepairBy != nil {
		repairBy, err = s.UserService.GetUserInfor(ctx, *repair.RepairBy)
		if err != nil {
			return nil, err
		}
	}

	location, err := s.LocationService.GetLocationByID(ctx, repair.Location)
	if err != nil {
		return nil, err
	}

	imageReport := make([]string, len(repair.ImageReport))
	for i, image := range repair.ImageReport {
		imageKey, err := s.UploaderService.GetImageKey(ctx, image)
		if err != nil {
			imageReport[i] = ""
		} else if imageKey == nil {
			imageReport[i] = ""
		} else {
			imageReport[i] = imageKey.Url
		}
	}

	var imageRepair []string
	if repair.ImageRepair != nil {
		imageRepair = make([]string, len(repair.ImageRepair))
		for i, image := range repair.ImageRepair {
			if image != nil {
				imageKey, err := s.UploaderService.GetImageKey(ctx, *image)
				if err != nil {
					imageRepair[i] = ""
				} else if imageKey == nil {
					imageRepair[i] = ""
				} else {
					imageRepair[i] = imageKey.Url
				}
			} else {
				imageRepair[i] = ""
			}
		}
	} else {
		log.Printf("GetRepairByID: repair.ImageRepair is nil")
	}

	// Get shop items for this repair
	shopItems, err := s.ShopService.GetRepairItems(ctx, repair.ID)
	if err != nil {
		log.Printf("GetRepairByID: Error getting shop items: %v", err)
		// Continue without shop items if error
		shopItems = nil
	}

	result := buildRepairResponse(repair, reportBy, repairBy, location, imageReport, imageRepair, shopItems)

	return result, nil
}

func (s *repairService) UpdateRepair(ctx context.Context, req UpdateRepairRequest, id string, userID string) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	existingRepair, err := s.RepairRepo.GetRepairByID(ctx, objectID)
	if err != nil {
		return err
	}

	if err := s.Policy.CanUpdateReport(ctx, existingRepair, userID); err != nil {
		return err
	}

	if req.JobName != "" {
		existingRepair.JobName = req.JobName
	}

	if req.Location != "" {
		existingRepair.Location = req.Location
	}

	if req.UrgentVote > 0 {
		existingRepair.UrgentVote = req.UrgentVote
	}

	if req.Comment != "" {
		existingRepair.Comment = req.Comment
	}

	if len(req.ImageReport) > 0 {
		for _, image := range existingRepair.ImageReport {
			err := s.UploaderService.DeleteImageKey(ctx, image)
			if err != nil {
				return err
			}
		}
		existingRepair.ImageReport = req.ImageReport
	}

	existingRepair.UpdatedAt = time.Now()

	err = s.RepairRepo.UpdateRepair(ctx, objectID, existingRepair)
	if err != nil {
		return err
	}

	return nil
}

func (s *repairService) DeleteRepair(ctx context.Context, id string, userID string) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	existingRepair, err := s.RepairRepo.GetRepairByID(ctx, objectID)
	if err != nil {
		return err
	}

	if err := s.Policy.CanDeleteRepair(ctx, existingRepair, userID); err != nil {
		return err
	}

	if len(existingRepair.ImageReport) > 0 {
		for _, image := range existingRepair.ImageReport {
			err := s.UploaderService.DeleteImageKey(ctx, image)
			if err != nil {
				return err
			}
		}
	}

	if len(existingRepair.ImageRepair) > 0 {
		for _, image := range existingRepair.ImageRepair {
			err := s.UploaderService.DeleteImageKey(ctx, *image)
			if err != nil {
				return err
			}
		}
	}

	err = s.RepairRepo.DeleteRepair(ctx, objectID)
	if err != nil {
		return err
	}

	return nil
}

func (s *repairService) AssignRepair(ctx context.Context, id string, userID string, req AssignRepairRequest) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	existingRepair, err := s.RepairRepo.GetRepairByID(ctx, objectID)
	if err != nil {
		return err
	}

	if err := s.Policy.CanAssignRepair(ctx, existingRepair, userID); err != nil {
		return err
	}

	existingRepair.AssignedTo = &req.AssignedTo
	existingRepair.Status = "assigned"
	existingRepair.UpdatedAt = time.Now()

	err = s.RepairRepo.UpdateRepair(ctx, objectID, existingRepair)
	if err != nil {
		return err
	}

	return nil

}

func (s *repairService) CompleteRepair(ctx context.Context, id string, req CompleteRepairRequest, userID string) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	existingRepair, err := s.RepairRepo.GetRepairByID(ctx, objectID)
	if err != nil {
		return err
	}

	if err := s.Policy.CanCompleteRepair(ctx, existingRepair, userID); err != nil {
		return err
	}

	if req.CommentRepair != nil {
		existingRepair.CommentRepair = req.CommentRepair
	}

	if len(req.ImageRepair) > 0 {
		// Convert []string to []*string for storage
		imageRepairPtrs := make([]*string, len(req.ImageRepair))
		for i, img := range req.ImageRepair {
			imageRepairPtrs[i] = &img
		}
		existingRepair.ImageRepair = imageRepairPtrs
	}

	// Calculate total cost from shop items
	shopItemsSummary, err := s.ShopService.GetRepairItems(ctx, objectID)
	if err == nil && shopItemsSummary != nil {
		existingRepair.TotalCost = &shopItemsSummary.TotalCost
	}

	now := time.Now()
	existingRepair.Status = "completed"
	existingRepair.DateRepair = &now
	existingRepair.RepairBy = &userID
	existingRepair.UpdatedAt = now

	err = s.RepairRepo.UpdateRepair(ctx, objectID, existingRepair)
	if err != nil {
		return err
	}

	return nil
}
