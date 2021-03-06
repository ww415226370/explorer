package controller

import (
	"strconv"

	"github.com/gorilla/mux"
	"github.com/irisnet/explorer/backend/model"
	"github.com/irisnet/explorer/backend/service"
	"github.com/irisnet/explorer/backend/types"
	"github.com/irisnet/explorer/backend/utils"
)

func RegisterProposal(r *mux.Router) error {
	funs := []func(*mux.Router) error{
		registerQueryProposals,
		registerQueryProposal,
		registerQueryGovParams,
		registerQueryDepositAndVotingProposals,
		registerQueryProposalDepositorTxs,
		registerQueryProposalVoterTxs,
	}

	for _, fn := range funs {
		if err := fn(r); err != nil {
			return err
		}
	}
	return nil
}

type Gov struct {
	*service.ProposalService
	*service.GovParamsService
}

var gov = Gov{
	service.Get(service.Proposal).(*service.ProposalService),
	service.Get(service.GovParams).(*service.GovParamsService),
}

func registerQueryProposals(r *mux.Router) error {

	doApi(r, types.UrlRegisterQueryProposals, "GET", func(request model.IrisReq) interface{} {
		page := int(utils.ParseIntWithDefault(QueryParam(request, "page"), 1))
		size := int(utils.ParseIntWithDefault(QueryParam(request, "size"), 10))

		result := gov.QueryList(page, size)
		return result
	})

	return nil
}

func registerQueryDepositAndVotingProposals(r *mux.Router) error {

	doApi(r, types.UrlRegisterQueryDepositVotingProposals, "GET", func(request model.IrisReq) interface{} {

		result := gov.QueryDepositAndVotingProposalList()
		return result
	})

	return nil
}

func registerQueryProposal(r *mux.Router) error {

	doApi(r, types.UrlRegisterQueryProposal, "GET", func(request model.IrisReq) interface{} {
		pid, err := strconv.Atoi(Var(request, "pid"))
		if err != nil {
			panic(types.CodeInValidParam)
		}

		result := gov.Query(pid)
		return result
	})

	return nil
}

func registerQueryGovParams(r *mux.Router) error {
	doApi(r, types.UrlRegisterQueryGovParams, "GET", func(request model.IrisReq) interface{} {
		return gov.GovParamsService.QueryAll()
	})
	return nil
}

func registerQueryProposalVoterTxs(r *mux.Router) error {
	doApi(r, types.UrlRegisterQueryProposalsVoterTxs, "GET", func(request model.IrisReq) interface{} {
		id, err := strconv.ParseInt(Var(request, "id"), 10, 64)
		if err != nil {
			panic(types.CodeInValidParam)
		}
		page := int(utils.ParseIntWithDefault(QueryParam(request, "page"), 1))
		size := int(utils.ParseIntWithDefault(QueryParam(request, "size"), 10))

		return gov.GetVoteTxs(id, page, size)
	})
	return nil
}

func registerQueryProposalDepositorTxs(r *mux.Router) error {
	doApi(r, types.UrlRegisterQueryProposalsDepositorTxs, "GET", func(request model.IrisReq) interface{} {
		id, err := strconv.ParseInt(Var(request, "id"), 10, 64)
		if err != nil {
			panic(types.CodeInValidParam)
		}
		page := int(utils.ParseIntWithDefault(QueryParam(request, "page"), 1))
		size := int(utils.ParseIntWithDefault(QueryParam(request, "size"), 10))

		return gov.GetDepositTxs(id, page, size)
	})
	return nil
}
