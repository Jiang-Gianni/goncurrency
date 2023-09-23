/*

Use context.WithValue when:
- the data should transit process or API boundaries
- the data should be immutable
- the data should trend toward simple types
- the data should be data, not types with methods
- the data should help decorate operations, not drive them

OK Examples: RequestID, UserID, Authorization Token

Not OK Example: Server Connection

*/

package chapt4

import (
	"context"
	"fmt"
)

type ctxKey int

const (
	ctxUserID ctxKey = iota
	ctxAuthToken
)

func ContextWithValue() {
	ProcessRequest("Jane", "abc123")
}

func ProcessRequest(userID string, authToken string) {
	ctx := context.WithValue(context.Background(), ctxUserID, userID)
	ctx = context.WithValue(ctx, ctxAuthToken, authToken)
	HandleResponse(ctx)
}

func HandleResponse(ctx context.Context) {
	fmt.Printf(
		"handling response for %v (%v)\n",
		ctx.Value(ctxUserID).(string),
		ctx.Value(ctxAuthToken).(string),
	)
}
