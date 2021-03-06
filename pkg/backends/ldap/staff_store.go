package ldap

import (
	"log"

	"github.com/go-ldap/ldap"

	"github.com/liut/staffio/pkg/models"
)

var (
	// simpleSecurityObject, "simpleSecurityObject"
	objectClassPeople = []string{"top", "staffioPerson", "uidObject", "inetOrgPerson"}
)

func (ls *ldapSource) StoreStaff(staff *models.Staff) (isNew bool, err error) {
	uid := staff.Uid
	err = ls.Bind(ls.BindDN, ls.Passwd, false)
	if err != nil {
		return
	}
	dn := ls.UDN(uid)
	var entry *ldap.Entry
	entry, err = ls.getEntry(ls.UDN(uid))
	if err == nil {
		// :update
		mr := makeModifyRequest(dn, entry, staff)
		if staff.EmployeeNumber != entry.GetAttributeValue("employeeNumber") {
			mr.Replace("employeeNumber", []string{staff.EmployeeNumber})
		}
		if staff.EmployeeType != entry.GetAttributeValue("employeeType") {
			mr.Replace("employeeType", []string{staff.EmployeeType})
		}
		err = ls.c.Modify(mr)
		if err != nil {
			log.Printf("modify %v ERR %s", mr, err)
		}
		return
	}
	if err == ErrNotFound {
		isNew = true
		ar := makeAddRequest(dn, staff)
		err = ls.c.Add(ar)
		if err != nil {
			log.Printf("add %v ERR %s", ar, err)
		}
		return
	}

	log.Printf("getEntry %s ERR %s", uid, err)

	return
}

func makeAddRequest(dn string, staff *models.Staff) *ldap.AddRequest {
	ar := ldap.NewAddRequest(dn)
	ar.Attribute("objectClass", objectClassPeople)
	ar.Attribute("uid", []string{staff.Uid})
	ar.Attribute("sn", []string{staff.Surname})
	ar.Attribute("givenName", []string{staff.GivenName})
	ar.Attribute("cn", []string{staff.GetCommonName()})
	ar.Attribute("mail", []string{staff.Email})
	if staff.Nickname != "" {
		ar.Attribute("displayName", []string{staff.Nickname})
	}
	if staff.Mobile != "" {
		ar.Attribute("mobile", []string{staff.Mobile})
	}

	if staff.EmployeeNumber != "" {
		ar.Attribute("employeeNumber", []string{staff.EmployeeNumber})
	}
	if staff.EmployeeType != "" {
		ar.Attribute("employeeType", []string{staff.EmployeeType})
	}
	if staff.Gender != models.Unknown {
		ar.Attribute("gender", []string{staff.Gender.String()[0:1]})
	}
	if staff.Birthday != "" {
		ar.Attribute("dateOfBirth", []string{staff.Birthday})
	}
	if staff.Description != "" {
		ar.Attribute("description", []string{staff.Description})
	}
	if staff.AvatarPath != "" {
		ar.Attribute("avatarPath", []string{staff.AvatarPath})
	}

	// if staff.Passwd != "" {
	// 	ar.Attribute("userPassword", []string{staff.Passwd})
	// }

	return ar
}

func makeModifyRequest(dn string, entry *ldap.Entry, staff *models.Staff) *ldap.ModifyRequest {
	mr := ldap.NewModifyRequest(dn)
	mr.Replace("objectClass", objectClassPeople)
	if staff.Surname != entry.GetAttributeValue("sn") {
		mr.Replace("sn", []string{staff.Surname})
	}
	if staff.GivenName != entry.GetAttributeValue("givenName") {
		mr.Replace("givenName", []string{staff.GivenName})
	}
	if staff.CommonName != entry.GetAttributeValue("cn") {
		mr.Replace("cn", []string{staff.GetCommonName()})
	}
	if staff.Nickname != "" && staff.Nickname != entry.GetAttributeValue("displayName") {
		mr.Replace("displayName", []string{staff.Nickname})
	}
	if staff.Email != entry.GetAttributeValue("mail") {
		mr.Replace("mail", []string{staff.Email})
	}
	if staff.Mobile != entry.GetAttributeValue("mobile") {
		mr.Replace("mobile", []string{staff.Mobile})
	}
	if staff.AvatarPath != entry.GetAttributeValue("avatarPath") {
		mr.Replace("avatarPath", []string{staff.AvatarPath})
	}
	if staff.Gender != models.Unknown {
		mr.Replace("gender", []string{staff.Gender.String()[0:1]})
	}
	if staff.Birthday != entry.GetAttributeValue("dateOfBirth") {
		mr.Replace("dateOfBirth", []string{staff.Birthday})
	}
	if staff.Description != entry.GetAttributeValue("description") {
		mr.Replace("description", []string{staff.Description})
	}
	return mr
}

/*
uid
sn
givenName
cn
mail
displayName
mobile
employeeNumber
employeeType
description
*/
