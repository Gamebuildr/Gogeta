package storehouse

import "testing"

type MockStorage struct {
	UploadCalled bool
}

type MockCompression struct {
	MockLocation string
	EncodeCalled bool
}

func (mock *MockStorage) Upload(data *StorageData) error {
	mock.UploadCalled = true
	return nil
}

func (mock *MockCompression) Encode(source string, target string) error {
	mock.EncodeCalled = true
	return nil
}

func TestCompressedStorehouseImplementsInterface(t *testing.T) {
	mockData := new(StorageData)
	compressStoreHouse := new(Compressed)
	mockCompression := MockCompression{
		MockLocation: "/mocked/Location",
		EncodeCalled: false,
	}
	mockStorage := MockStorage{UploadCalled: false}

	compressStoreHouse.Compression = &mockCompression
	compressStoreHouse.StorageSystem = &mockStorage

	compressStoreHouse.StoreFiles(mockData)

	if !mockCompression.EncodeCalled {
		t.Errorf("Expected: %v", "Encode to be called")
	}
	if !mockStorage.UploadCalled {
		t.Errorf("Expected: %v", "Upload to be called")
	}
}
