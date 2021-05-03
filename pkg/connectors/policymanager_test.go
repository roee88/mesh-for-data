package connectors_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mesh-for-data/mesh-for-data/pkg/connectors"
	pb "github.com/mesh-for-data/mesh-for-data/pkg/connectors/protobuf"
)

var _ = Describe("PolicyManager", func() {
	Describe("merge policy decisions", func() {

		removeColumn1 := &pb.EnforcementAction{Name: "remove column", Id: "remove-ID", Level: pb.EnforcementAction_COLUMN, Args: map[string]string{"column_name": "col1"}}
		removeColumn2 := &pb.EnforcementAction{Name: "remove column", Id: "remove-ID", Level: pb.EnforcementAction_COLUMN, Args: map[string]string{"column_name": "col2"}}
		redactColumn1 := &pb.EnforcementAction{Name: "redact column", Id: "redact-ID", Level: pb.EnforcementAction_COLUMN, Args: map[string]string{"column_name": "col1"}}

		Context("on same dataset", func() {

			Context("with same operation (read)", func() {

				Context("has multiple actions on the same column", func() {

					It("should put the actions in the same EnforcementActions slice", func() {
						left := &pb.PoliciesDecisions{DatasetDecisions: []*pb.DatasetDecision{{
							Dataset: &pb.DatasetIdentifier{DatasetId: "1"},
							Decisions: []*pb.OperationDecision{
								{
									Operation:          &pb.AccessOperation{Type: pb.AccessOperation_READ},
									EnforcementActions: []*pb.EnforcementAction{removeColumn1},
								},
							},
						}}}
						right := &pb.PoliciesDecisions{DatasetDecisions: []*pb.DatasetDecision{{
							Dataset: &pb.DatasetIdentifier{DatasetId: "1"},
							Decisions: []*pb.OperationDecision{
								{
									Operation:          &pb.AccessOperation{Type: pb.AccessOperation_READ},
									EnforcementActions: []*pb.EnforcementAction{redactColumn1},
								},
							},
						}}}
						expected := &pb.PoliciesDecisions{DatasetDecisions: []*pb.DatasetDecision{{
							Dataset: &pb.DatasetIdentifier{DatasetId: "1"},
							Decisions: []*pb.OperationDecision{
								{
									Operation:          &pb.AccessOperation{Type: pb.AccessOperation_READ},
									EnforcementActions: []*pb.EnforcementAction{removeColumn1, redactColumn1},
								},
							},
						}}}
						Expect(connectors.MergePoliciesDecisions(left, right)).To(Equal(expected))
					})
				})

				Context("has same action type but on different columns", func() {
					It("put the actions in the same EnforcementActions slice", func() {
						left := &pb.PoliciesDecisions{DatasetDecisions: []*pb.DatasetDecision{{
							Dataset: &pb.DatasetIdentifier{DatasetId: "1"},
							Decisions: []*pb.OperationDecision{
								{
									Operation:          &pb.AccessOperation{Type: pb.AccessOperation_READ},
									EnforcementActions: []*pb.EnforcementAction{removeColumn1},
								},
							},
						}}}
						right := &pb.PoliciesDecisions{DatasetDecisions: []*pb.DatasetDecision{{
							Dataset: &pb.DatasetIdentifier{DatasetId: "1"},
							Decisions: []*pb.OperationDecision{
								{
									Operation:          &pb.AccessOperation{Type: pb.AccessOperation_READ},
									EnforcementActions: []*pb.EnforcementAction{removeColumn2},
								},
							},
						}}}
						expected := &pb.PoliciesDecisions{DatasetDecisions: []*pb.DatasetDecision{{
							Dataset: &pb.DatasetIdentifier{DatasetId: "1"},
							Decisions: []*pb.OperationDecision{
								{
									Operation:          &pb.AccessOperation{Type: pb.AccessOperation_READ},
									EnforcementActions: []*pb.EnforcementAction{removeColumn1, removeColumn2},
								},
							},
						}}}
						Expect(connectors.MergePoliciesDecisions(left, right)).To(Equal(expected))
					})
				})
			})

			Context("with multiple operations (read, write)", func() {

				Context("has same action on the same column", func() {
					It("should result in two decisions for the dataset", func() {
						left := &pb.PoliciesDecisions{DatasetDecisions: []*pb.DatasetDecision{{
							Dataset: &pb.DatasetIdentifier{DatasetId: "1"},
							Decisions: []*pb.OperationDecision{
								{
									Operation:          &pb.AccessOperation{Type: pb.AccessOperation_READ},
									EnforcementActions: []*pb.EnforcementAction{removeColumn1},
								},
							},
						}}}
						right := &pb.PoliciesDecisions{DatasetDecisions: []*pb.DatasetDecision{{
							Dataset: &pb.DatasetIdentifier{DatasetId: "1"},
							Decisions: []*pb.OperationDecision{
								{
									Operation:          &pb.AccessOperation{Type: pb.AccessOperation_WRITE},
									EnforcementActions: []*pb.EnforcementAction{removeColumn1},
								},
							},
						}}}
						expected := &pb.PoliciesDecisions{DatasetDecisions: []*pb.DatasetDecision{{
							Dataset: &pb.DatasetIdentifier{DatasetId: "1"},
							Decisions: []*pb.OperationDecision{
								{
									Operation:          &pb.AccessOperation{Type: pb.AccessOperation_READ},
									EnforcementActions: []*pb.EnforcementAction{removeColumn1},
								},
								{
									Operation:          &pb.AccessOperation{Type: pb.AccessOperation_WRITE},
									EnforcementActions: []*pb.EnforcementAction{removeColumn1},
								},
							},
						}}}
						Expect(connectors.MergePoliciesDecisions(left, right)).To(Equal(expected))
					})
				})

			})
		})

		Context("on two datasets", func() {
			It("should keep as separate dataset decisions", func() {
				left := &pb.PoliciesDecisions{DatasetDecisions: []*pb.DatasetDecision{{
					Dataset: &pb.DatasetIdentifier{DatasetId: "1"},
					Decisions: []*pb.OperationDecision{
						{
							Operation:          &pb.AccessOperation{Type: pb.AccessOperation_READ},
							EnforcementActions: []*pb.EnforcementAction{removeColumn1},
						},
					},
				}}}
				right := &pb.PoliciesDecisions{DatasetDecisions: []*pb.DatasetDecision{{
					Dataset: &pb.DatasetIdentifier{DatasetId: "2"},
					Decisions: []*pb.OperationDecision{
						{
							Operation:          &pb.AccessOperation{Type: pb.AccessOperation_READ},
							EnforcementActions: []*pb.EnforcementAction{removeColumn1},
						},
					},
				}}}
				expected := &pb.PoliciesDecisions{DatasetDecisions: []*pb.DatasetDecision{
					{
						Dataset: &pb.DatasetIdentifier{DatasetId: "1"},
						Decisions: []*pb.OperationDecision{
							{
								Operation:          &pb.AccessOperation{Type: pb.AccessOperation_READ},
								EnforcementActions: []*pb.EnforcementAction{removeColumn1},
							},
						},
					},
					{
						Dataset: &pb.DatasetIdentifier{DatasetId: "2"},
						Decisions: []*pb.OperationDecision{
							{
								Operation:          &pb.AccessOperation{Type: pb.AccessOperation_READ},
								EnforcementActions: []*pb.EnforcementAction{removeColumn1},
							},
						},
					},
				}}
				Expect(connectors.MergePoliciesDecisions(left, right)).To(Equal(expected))
			})
		})

	})

})
