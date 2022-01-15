# OptionDance MatchEngine

Option Dance engine is based on the Ocean.one match engine(https://github.com/MixinNetwork/ocean.one#bid-order-behavior)

## OrderBook V1 data structure

As ocean.one engine with it bid-order-behavior (https://github.com/MixinNetwork/ocean.one#bid-order-behavior)


An bid order would be like: 
```go
// not filled
bid := &core.EngineOrder{
    Id:              id.String(),
    Side:            PageSideBid,
    Type:            core.OrderTypeLimit,
    Price:           "100",
    RemainingAmount: "0",
    FilledAmount:    "0",
    RemainingFunds:  "100",
    FilledFunds:     "0",
}

// paritical filled
bid := &core.EngineOrder{
    Id:              id.String(),
    Side:            PageSideBid,
    Type:            core.OrderTypeLimit,
    Price:           "100",
    RemainingAmount: "0",
    FilledAmount:    "0.5",
    RemainingFunds:  "50",
    FilledFunds:     "50",
}

//full filled
ask := &core.EngineOrder{
    Id:              id.String(),
    Side:            PageSideAsk,
    Type:            core.OrderTypeLimit,
    Price:           "100",
    RemainingAmount: "0",
    FilledAmount:    "1",
    RemainingFunds:  "0",
    FilledFunds:     "100",
}
```

note that in bid order, only the Price and RemainingFunds is positive , the remainingAmount,filledAmount,filledFunds is all equals zero

An ask order will be like:
```go

// not filled
ask := &core.EngineOrder{
    Id:              id.String(),
    Side:            PageSideAsk,
    Type:            core.OrderTypeLimit,
    Price:           "100",
    RemainingAmount: "1",
    FilledAmount:    "0",
    RemainingFunds:  "0",
    FilledFunds:     "0",
}

// partially filled
ask := &core.EngineOrder{
    Id:              id.String(),
    Side:            PageSideAsk,
    Type:            core.OrderTypeLimit,
    Price:           "100",
    RemainingAmount: "0.5",
    FilledAmount:    "0.5",
    RemainingFunds:  "0",
    FilledFunds:     "50",
}

//full filled
ask := &core.EngineOrder{
    Id:              id.String(),
    Side:            PageSideAsk,
    Type:            core.OrderTypeLimit,
    Price:           "100",
    RemainingAmount: "0",
    FilledAmount:    "1",
    RemainingFunds:  "0",
    FilledFunds:     "100",
}
```
note that in bid order, only the Price and RemainingAmount is positive, the RemainingFunds,filledAmount,filledFunds is all equals zero


## OrderBook V2 data structure

In order to support the case that bid order should be finished when remainingAmount is zero, then the bid order amount-funds status should be
adjusted like the following

```go

// not filled
bid := &core.EngineOrder{
    Id:              id.String(),
    Side:            PageSideBid,
    Type:            core.OrderTypeLimit,
    Price:           "100",
    RemainingAmount: "1",
    FilledAmount:    "0",
    RemainingFunds:  "100",
    FilledFunds:     "0",
}

// partially filled
bid := &core.EngineOrder{
    Id:              id.String(),
    Side:            PageSideBid,
    Type:            core.OrderTypeLimit,
    Price:           "100",
    RemainingAmount: "0.5",
    FilledAmount:    "0.5",
    RemainingFunds:  "50",
    FilledFunds:     "50",
}

//full filled
ask := &core.EngineOrder{
    Id:              id.String(),
    Side:            PageSideAsk,
    Type:            core.OrderTypeLimit,
    Price:           "100",
    RemainingAmount: "0",
    FilledAmount:    "1",
    RemainingFunds:  "0",
    FilledFunds:     "100",
}

```