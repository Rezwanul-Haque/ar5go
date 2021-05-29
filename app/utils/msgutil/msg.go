package msgutil

import "fmt"

// Success Massage
func EntityUploadSuccessMsg(entity string) string {
	return fmt.Sprintf("%s uploaded successfully", entity)
}

func EntityGetSuccessMsg(entity string) string {
	return fmt.Sprintf("%s get successful", entity)
}

func EntityGetEmptySuccessMsg(entity string) string {
	return fmt.Sprintf("%s has no data", entity)
}

func EntityUpdateSuccessMsg(entity string) string {
	return fmt.Sprintf("%s updated successfully", entity)
}

func EntityChangedSuccessMsg(entity string) string {
	return fmt.Sprintf("%s changed successfully", entity)
}

func EntityDeleteSuccessMsg(entity string) string {
	return fmt.Sprintf("%s deleted successfully", entity)
}

func EntityNotFoundMsg(entity string) string {
	return fmt.Sprintf("%s not found", entity)
}

//Failed Message
func EntityCreationFailedMsg(entity string) string {
	return fmt.Sprintf("failed to create %s", entity)
}

func EntityStructToStructFailedMsg(entity string) string {
	return fmt.Sprintf("error occur when unmarshalling from struct to struct - %s", entity)
}

func EntityBindToStructFailedMsg(entity string) string {
	return fmt.Sprintf("error occur when bind from request to struct - %s", entity)
}

// Generic failed message
func SomethingWentWrongMsg() string {
	return "something went wrong"
}

func EntityGenericFailedMsg(entity string) string {
	return fmt.Sprintf("error occur when getting %s", entity)
}

func EntityGenericInvalidMsg(entity string) string {
	return fmt.Sprintf("invalid %s", entity)
}
