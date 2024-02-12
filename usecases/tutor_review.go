package usecases

import (
	"github.com/Arun4rangan/api-rfrl/rfrl"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
)

type TutorReviewUseCase struct {
	DB               rfrl.DB
	TutorReviewStore rfrl.TutorReviewStore
	SessionStore     rfrl.SessionStore
	ClientStore      rfrl.ClientStore
}

func NewTutorReviewUseCase(
	db rfrl.DB,
	tutorReviewStore rfrl.TutorReviewStore,
	sessionStore rfrl.SessionStore,
	clientStore rfrl.ClientStore,
) *TutorReviewUseCase {
	return &TutorReviewUseCase{db, tutorReviewStore, sessionStore, clientStore}
}

func (tru *TutorReviewUseCase) CreateTutorReview(ClientID string, TutorID string, Stars int, Review string, Headline string) (*rfrl.TutorReview, error) {
	if ClientID == TutorID {
		return nil, errors.Errorf("Client cannot create a review for said client")
	}

	checkOverlap, err := tru.SessionStore.CheckClientsAttendedTutorSession(tru.DB, TutorID, []string{ClientID})

	if err != nil {
		return nil, err
	}

	if checkOverlap == false {
		return nil, errors.Errorf("Client did not get tutored by Tutor")
	}

	tutorReview := rfrl.NewTutorReview(TutorID, Stars, Review, Headline)

	createdTutorReview, err := tru.TutorReviewStore.CreateTutorReview(tru.DB, ClientID, &tutorReview)

	if err != nil {
		return nil, err
	}

	err = tru.TutorReviewStore.DeletePendingReview(tru.DB, ClientID, TutorID)

	if err != nil {
		return nil, err
	}

	fromClient, err := tru.ClientStore.GetClientFromID(tru.DB, createdTutorReview.FromID)

	if err != nil {
		return nil, err
	}

	createdTutorReview.FromClient = *fromClient

	return createdTutorReview, err
}

func (tru *TutorReviewUseCase) UpdateTutorReview(ClientID string, ID int, Stars int, Review string, Headline string) (*rfrl.TutorReview, error) {
	forClient, err := tru.TutorReviewStore.CheckTutorReviewForClient(tru.DB, ClientID, ID)

	if err != nil {
		return nil, err
	}

	if forClient == false {
		return nil, errors.Errorf("Tutor Review  (%d) does not belong to this client %s", ID, ClientID)
	}

	tutorReview := rfrl.TutorReview{}
	tutorReview.ID = ID
	tutorReview.Stars = null.NewInt(int64(Stars), true)
	tutorReview.Review = null.NewString(Review, true)
	tutorReview.Headline = null.NewString(Headline, true)

	return tru.TutorReviewStore.UpdateTutorReview(tru.DB, &tutorReview)
}

func (tru *TutorReviewUseCase) DeleteTutorReview(ClientID string, ID int) error {
	forClient, err := tru.TutorReviewStore.CheckTutorReviewForClient(tru.DB, ClientID, ID)

	if err != nil {
		return err
	}

	if forClient == false {
		return errors.Errorf("Tutor Review  (%d) does not belong to this client %s", ID, ClientID)
	}

	return tru.TutorReviewStore.DeleteTutorReview(tru.DB, ID)
}

func (tru *TutorReviewUseCase) GetTutorReview(ID int) (*rfrl.TutorReview, error) {
	tutorReview, err := tru.TutorReviewStore.GetTutorReview(tru.DB, ID)
	if err != nil {
		return nil, err
	}
	fromClient, err := tru.ClientStore.GetClientFromID(tru.DB, tutorReview.FromID)

	if err != nil {
		return nil, err
	}

	tutorReview.FromClient = *fromClient

	return tutorReview, nil
}

func (tru *TutorReviewUseCase) GetTutorReviews(TutorID string) (*[]rfrl.TutorReview, error) {
	tutorReviews, err := tru.TutorReviewStore.GetTutorReviews(tru.DB, TutorID)
	if err != nil {
		return nil, err
	}

	clientIDToIndex := make(map[string]int)
	var clientIDs []string
	for i := 0; i < len(*tutorReviews); i++ {
		clientIDs = append(clientIDs, (*tutorReviews)[i].FromID)
		clientIDToIndex[(*tutorReviews)[i].FromID] = i
	}

	if len(clientIDs) > 0 {
		clients, err := tru.ClientStore.GetClientFromIDs(tru.DB, clientIDs)

		if err != nil {
			return nil, err
		}

		for i := 0; i < len(*clients); i++ {
			index := clientIDToIndex[(*clients)[i].ID]
			(*tutorReviews)[index].FromClient = (*clients)[i]
		}
	}

	return tutorReviews, nil
}

func (tru *TutorReviewUseCase) GetTutorReviewsAggregate(ClientID string) (*rfrl.TutorReviewAggregate, error) {
	return tru.TutorReviewStore.GetTutorReviewsAggregate(tru.DB, ClientID)
}

func (tru *TutorReviewUseCase) GetPendingReviews(ClientID string) (*[]rfrl.PendingTutorReview, error) {
	return tru.TutorReviewStore.GetPendingReviews(tru.DB, ClientID)
}

func (tru *TutorReviewUseCase) CreatePendingReview(menteeID string, tutorID string) error {
	alreadyExist, err := tru.TutorReviewStore.CheckIfReviewAlreadyExists(tru.DB, menteeID, tutorID)

	if err != nil {
		return err
	}

	if alreadyExist {
		return nil
	}

	return tru.TutorReviewStore.CreatePendingReview(tru.DB, menteeID, tutorID)
}

func (tru *TutorReviewUseCase) DeletePendingReview(menteeID string, tutorID string) error {
	return tru.TutorReviewStore.DeletePendingReview(tru.DB, menteeID, tutorID)
}
