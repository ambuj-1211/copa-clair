package main

import (
	"encoding/json"
	"fmt"
	"os"

	wizTypes "github.com/ambuj-1211/copa-wiz" //TODO : GET THE REAL TYPE PACKAGE FOR WIZ
	v1alpha1 "github.com/project-copacetic/copacetic/pkg/types/v1alpha1"
)

type WizParser struct{}

// parseFakeReport parses a fake report from a file
func parseWizReport(file string) (*wizTypes.Document, error) {// TODO : Update according to the wiz types.
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var wr wizTypes.Document // TODO: Update according to the wiz types
	if err = json.Unmarshal(data, &wr); err != nil {
		return nil, err
	}

	return &wr, nil
}

func NewWizParser() *WizParser {
	return &WizParser{}
}

func (k *WizParser) Parse(file string) (*v1alpha1.UpdateManifest, error) {
	// Parse the fake report
	report, err := parseWizReport(file)
	if err != nil {
		return nil, err
	}

	if report.Descriptor.Name != "wiz" {
		return nil, errors.New("report format not supported by wiz")
	}

	if err != nil {
		return nil, err
	}

	// Create the standardized report
	// TODO: Use the report variable according to what fetched for wiz.
	updates := v1alpha1.UpdateManifest{
		APIVersion: v1alpha1.APIVersion,
		Metadata: v1alpha1.Metadata{
			OS: v1alpha1.OS{
				Type: report.Distro.Name,
				Version: report.Distro.Version,
			},
			Config: v1alpha1.Config{
				Arch: report.Source.Target.(map[string]interface{})["architecture"].(string),
			},
		},
	}
	for i := range report.Matches {
		vuln := &report.Matches[i]
		if vuln.Artifact.Language == "" && vuln.Vulnerability.Fix.State == "fixed" {
			updates.Updates = append(updates.Updates, v1alpha1.UpdatePackage{Name: vuln.Artifact.Name, InstalledVersion: vuln.Artifact.Version, FixedVersion: vuln.Vulnerability.Fix.Versions[0], VulnerabilityID: vuln.Vulnerability.ID})
		}
	}
	return &updates, nil
	// Convert the fake report to the standardized report
	// for i := range report.Packages {
	// 	pkgs := &report.Packages[i]
	// 	if pkgs.FixedVersion != "" {
	// 		updates.Updates = append(updates.Updates, v1alpha1.UpdatePackage{
	// 			Name: pkgs.Name,
	// 			InstalledVersion: pkgs.InstalledVersion,
	// 			FixedVersion: pkgs.FixedVersion,
	// 			VulnerabilityID: pkgs.VulnerabilityID,
	// 		})
	// 	}
	// }
	// return &updates, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <image report>\n", os.Args[0])
		os.Exit(1)
	}

	// Initialize the parser
	wizParser := NewWizParser()

	// Get the image report from command line
	imageReport := os.Args[1]

	report, err := wizParser.parse(imageReport)
	if err != nil {
		fmt.Printf("error parsing report: %v\n", err)
		os.Exit(1)
	}

	// Serialize the standardized report and print it to stdout
	reportBytes, err := json.Marshal(report)
	if err != nil {
		fmt.Printf("Error serializing report: %v\n", err)
		os.Exit(1)
	}

	os.Stdout.Write(reportBytes)
}