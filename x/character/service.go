package character

import (
    "log"
    "encoding/json"
    "github.com/totegamma/concurrent/x/util"
    "github.com/totegamma/concurrent/x/core"
)

// Service is service of characters
type Service struct {
    repo *Repository
}

// NewService is for wire.go
func NewService(repo *Repository) *Service {
    return &Service{repo: repo}
}

// GetCharacters returns characters by owner and schema
func (s* Service) GetCharacters(owner string, schema string) ([]core.Character, error) {
    characters, err := s.repo.Get(owner, schema)
    if err != nil {
        log.Printf("error occured while GetCharacters in characterRepository. error: %v\n", err)
        return []core.Character{}, err
    }
    return characters, nil
}

// PutCharacter creates new character if the signature is valid
func (s* Service) PutCharacter(objectStr string, signature string, id string) error {

    var object signedObject
    err := json.Unmarshal([]byte(objectStr), &object)
    if err != nil {
        return err
    }

    if err := util.VerifySignature(objectStr, object.Signer, signature); err != nil {
        log.Println("verify signature err: ", err)
        return err
    }

    character := core.Character {
        ID: id,
        Author: object.Signer,
        Schema: object.Schema,
        Payload: objectStr,
        Signature: signature,
    }

    err = s.repo.Upsert(character)
    if err != nil {
        return err
    }

    return nil
}
