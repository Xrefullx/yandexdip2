package services

import (
	"github.com/Xrefullx/YanDip/client/model"
	"github.com/Xrefullx/YanDip/client/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

var (
	testFile = "sync_test_case.json"
)

type syncTest struct {
	name   string
	ext    map[uuid.UUID]int
	loc    []model.SecretMeta
	reqErr require.ErrorAssertionFunc
	result []SyncTask
}

func TestSync_CalcSync(t *testing.T) {
	secretID := uuid.New()
	nilUUID := uuid.Nil
	locID := rand.Int63n(100) + 1
	ver := rand.Intn(100) + 1

	timeStamp := pkg.MakeTimestamp()

	tests := []syncTest{
		//  new secretID == nil
		{
			name:   "new status new - upload",
			loc:    []model.SecretMeta{{SecretID: nilUUID, ID: locID, SecretVer: 1, StatusID: model.SecretStatuses["NEW"], TimeStamp: timeStamp}},
			result: []SyncTask{taskUploadNew(model.SecretMeta{ID: locID, SecretVer: 1, TimeStamp: timeStamp})},
			reqErr: require.NoError,
		},
		{name: "new status edited - upload",
			loc:    []model.SecretMeta{{SecretID: nilUUID, ID: locID, SecretVer: 1, StatusID: model.SecretStatuses["EDITED"], TimeStamp: timeStamp}},
			result: []SyncTask{taskUploadNew(model.SecretMeta{ID: locID, SecretVer: 1, TimeStamp: timeStamp})},

			reqErr: require.NoError,
		},
		{
			name:   "new status deleted - delete hard",
			loc:    []model.SecretMeta{{SecretID: nilUUID, ID: locID, SecretVer: 1, StatusID: model.SecretStatuses["DELETED"], TimeStamp: timeStamp}},
			result: []SyncTask{taskDeleteLocally(model.SecretMeta{ID: locID, TimeStamp: timeStamp})},
			reqErr: require.NoError,
		},

		//  exist secretID != nil, deleted not in ext list
		{
			name:   "exist edited/deleted - delete hard",
			loc:    []model.SecretMeta{{SecretID: secretID, ID: locID, SecretVer: ver, StatusID: model.SecretStatuses["EDITED"], TimeStamp: timeStamp}},
			result: []SyncTask{taskDeleteLocally(model.SecretMeta{ID: locID, TimeStamp: timeStamp})},
			reqErr: require.NoError,
		},
		{
			name:   "exist deleted/deleted - delete hard",
			loc:    []model.SecretMeta{{SecretID: secretID, ID: locID, SecretVer: ver, StatusID: model.SecretStatuses["DELETED"], TimeStamp: timeStamp}},
			result: []SyncTask{taskDeleteLocally(model.SecretMeta{ID: locID, TimeStamp: timeStamp})},
			reqErr: require.NoError,
		},
		{
			name:   "exist actual/deleted - delete hard",
			loc:    []model.SecretMeta{{SecretID: secretID, ID: locID, SecretVer: ver, StatusID: model.SecretStatuses["ACTUAL"], TimeStamp: timeStamp}},
			result: []SyncTask{taskDeleteLocally(model.SecretMeta{ID: locID, TimeStamp: timeStamp})},
			reqErr: require.NoError,
		},

		//  exist secretID != nil, ext no changes ver.loc == ver.ext
		{
			name:   "exist edited/no changes - upload",
			loc:    []model.SecretMeta{{SecretID: secretID, ID: locID, SecretVer: ver, StatusID: model.SecretStatuses["EDITED"], TimeStamp: timeStamp}},
			ext:    map[uuid.UUID]int{secretID: ver},
			result: []SyncTask{taskUpload(model.SecretMeta{SecretID: secretID, SecretVer: ver, TimeStamp: timeStamp})},
			reqErr: require.NoError,
		},
		{
			name:   "exist deleted/no changes - send delete",
			loc:    []model.SecretMeta{{SecretID: secretID, ID: locID, SecretVer: ver, StatusID: model.SecretStatuses["DELETED"], TimeStamp: timeStamp}},
			ext:    map[uuid.UUID]int{secretID: ver},
			result: []SyncTask{taskDeleteRemote(model.SecretMeta{SecretID: secretID, TimeStamp: timeStamp})},
			reqErr: require.NoError,
		},

		{
			name:   "exist actual/no changes - nil",
			loc:    []model.SecretMeta{{SecretID: secretID, ID: locID, SecretVer: ver, StatusID: model.SecretStatuses["ACTUAL"], TimeStamp: timeStamp}},
			ext:    map[uuid.UUID]int{secretID: ver},
			result: []SyncTask{},
			reqErr: require.NoError,
		},

		//  exist secretID != nil, ext changed ver.loc < ver.ext
		{
			name:   "exist edited/ changed - collision",
			loc:    []model.SecretMeta{{SecretID: secretID, ID: locID, SecretVer: ver, StatusID: model.SecretStatuses["EDITED"], TimeStamp: timeStamp}},
			ext:    map[uuid.UUID]int{secretID: ver + 1},
			result: []SyncTask{taskCollision(model.SecretMeta{ID: locID, SecretID: secretID, TimeStamp: timeStamp}, ver+1)},
			reqErr: require.NoError,
		},

		{
			name:   "exist deleted/ changed - send delete",
			loc:    []model.SecretMeta{{SecretID: secretID, ID: locID, SecretVer: ver, StatusID: model.SecretStatuses["DELETED"], TimeStamp: timeStamp}},
			ext:    map[uuid.UUID]int{secretID: ver + 1},
			result: []SyncTask{taskDeleteRemote(model.SecretMeta{SecretID: secretID, TimeStamp: timeStamp})},
			reqErr: require.NoError,
		},

		{
			name:   "exist actual/ changed - download",
			loc:    []model.SecretMeta{{SecretID: secretID, ID: locID, SecretVer: ver, StatusID: model.SecretStatuses["ACTUAL"], TimeStamp: timeStamp}},
			ext:    map[uuid.UUID]int{secretID: ver + 1},
			result: []SyncTask{taskDownload(model.SecretMeta{SecretID: secretID, TimeStamp: timeStamp})},
			reqErr: require.NoError,
		},
		{
			name:   "not exist/ exist - download",
			loc:    []model.SecretMeta{},
			ext:    map[uuid.UUID]int{secretID: ver},
			result: []SyncTask{taskDownloadNew(secretID)},
			reqErr: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			syncResult, err := SyncService{}.CalcSyncBatch(tt.ext, tt.loc)

			tt.reqErr(t, err)
			require.Equal(t, tt.result, syncResult)
		})

	}
}
