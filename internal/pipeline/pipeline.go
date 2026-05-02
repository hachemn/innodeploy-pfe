package pipeline

import (
	"log"
	"os/exec"
)

type Step struct {
	Name    string
	Command string
}
func RunPipeline(repo string) {
	log.Println("🚀 Pipeline started for:", repo)

	steps := []Step{
		{"Cleanup", "rm -rf project"},
		{"Clone", "git clone " + repo + " project"},
		{"Build Docker Image", "cd project && docker build -t innodeploy-app ."},
		{"Tag Image", "docker tag innodeploy-app nourhenhachem/innodeploy-app:latest"},
		{"Push Image", "docker push nourhenhachem/innodeploy-app:latest"},
		{"Deploy to Kubernetes", "cd project && kubectl apply -f k8s/deployment.yaml"},
	}

	for _, step := range steps {
		runStep(step)
	}

	log.Println("✅ Pipeline finished")
}
func runStep(step Step) {
	log.Println("➡️ Step:", step.Name)

	command := exec.Command("sh", "-c", step.Command)
	output, err := command.CombinedOutput()

	if err != nil {
		log.Println("❌ Error in step:", step.Name)
		log.Println(string(output))
		return
	}

	log.Println("✅ Success:", step.Name)
	log.Println(string(output))
}