package internalgrpc

import (
	"context"
	"time"

	"github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/server/grpc/gen"
	"github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/storage"
)

type eventService struct {
	gen.UnimplementedEventServiceServer

	app *app.App
}

const dateLayout = time.DateTime

func (e *eventService) CreateEvent(
	ctx context.Context,
	rqParams *gen.CreateEventRequest,
) (*gen.CreateEventReply, error) {
	startDt, err := time.Parse(dateLayout, rqParams.StartDt)
	if err != nil {
		return nil, err
	}

	endDt, err := time.Parse(dateLayout, rqParams.EndDt)
	if err != nil {
		return nil, err
	}

	notifyBefore, err := time.ParseDuration(rqParams.NotifyBefore)
	if err != nil {
		return nil, err
	}

	event, err := e.app.CreateEvent(
		ctx,
		rqParams.UserId,
		rqParams.Title,
		startDt,
		endDt,
		rqParams.Description,
		notifyBefore,
	)
	if err != nil {
		return nil, err
	}

	return &gen.CreateEventReply{
		EventId: event.ID,
	}, nil
}

func (e *eventService) UpdateEvent(ctx context.Context, rqEvent *gen.Event) (*gen.UpdateEventReply, error) {
	startDt, err := time.Parse(dateLayout, rqEvent.StartDt)
	if err != nil {
		return nil, err
	}

	endDt, err := time.Parse(dateLayout, rqEvent.EndDt)
	if err != nil {
		return nil, err
	}

	notifyBefore, err := time.ParseDuration(rqEvent.NotifyBefore)
	if err != nil {
		return nil, err
	}

	_, updErr := e.app.UpdateEvent(
		ctx,
		rqEvent.Id,
		rqEvent.UserId,
		rqEvent.Title,
		startDt,
		endDt,
		rqEvent.Description,
		notifyBefore,
	)
	if updErr != nil {
		return &gen.UpdateEventReply{Ok: false}, updErr
	}

	return &gen.UpdateEventReply{Ok: true}, err
}

func (e *eventService) DeleteEvent(
	ctx context.Context,
	deleteRq *gen.DeleteEventRequest,
) (*gen.DeleteEventReply, error) {
	err := e.app.DeleteEvent(ctx, deleteRq.EventId)
	if err != nil {
		return &gen.DeleteEventReply{Ok: false}, err
	}

	return &gen.DeleteEventReply{Ok: true}, nil
}

func (e *eventService) ListEventsForDay(
	ctx context.Context,
	listRq *gen.ListEventsRequest,
) (*gen.ListEventsReply, error) {
	startDt, err := time.Parse(time.DateOnly, listRq.StartDt)
	if err != nil {
		return nil, err
	}

	events, err := e.app.ListEventsForDay(ctx, listRq.UserId, startDt)
	if err != nil {
		return nil, err
	}

	resEvents := make([]*gen.Event, len(events))
	for ind, evt := range events {
		resEvents[ind] = dumpAppEvent(evt)
	}

	return &gen.ListEventsReply{Events: resEvents}, nil
}

func (e *eventService) ListEventsForWeek(
	ctx context.Context,
	listRq *gen.ListEventsRequest,
) (*gen.ListEventsReply, error) {
	startDt, err := time.Parse(time.DateOnly, listRq.StartDt)
	if err != nil {
		return nil, err
	}

	events, err := e.app.ListEventsForWeek(ctx, listRq.UserId, startDt)
	if err != nil {
		return nil, err
	}

	resEvents := make([]*gen.Event, len(events))
	for ind, evt := range events {
		resEvents[ind] = dumpAppEvent(evt)
	}

	return &gen.ListEventsReply{Events: resEvents}, nil
}

func (e *eventService) ListEventsForMonth(
	ctx context.Context,
	listRq *gen.ListEventsRequest,
) (*gen.ListEventsReply, error) {
	startDt, err := time.Parse(time.DateOnly, listRq.StartDt)
	if err != nil {
		return nil, err
	}

	events, err := e.app.ListEventsForMonth(ctx, listRq.UserId, startDt)
	if err != nil {
		return nil, err
	}

	resEvents := make([]*gen.Event, len(events))
	for ind, evt := range events {
		resEvents[ind] = dumpAppEvent(evt)
	}

	return &gen.ListEventsReply{Events: resEvents}, nil
}

func dumpAppEvent(event *storage.Event) *gen.Event {
	return &gen.Event{
		Id:           event.ID,
		Title:        event.Title,
		StartDt:      event.StartDT.Format(dateLayout),
		EndDt:        event.EndDT.Format(dateLayout),
		Description:  event.Description,
		UserId:       event.UserID,
		NotifyBefore: event.NotifyBefore.String(),
	}
}
