package pipeline

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

type Step struct {
	Name    string
	Command string
}

func RunPipeline(repo string) {
	log.Println("🚀 Pipeline started for:", repo)

	// 🟢 version unique
	version := fmt.Sprintf("v%d", time.Now().Unix())

	imageName := "nourhenhachem/innodeploy-app:" + version

	steps := []Step{
		{"Cleanup", "rm -rf project"},
		{"Clone", "git clone " + repo + " project"},

		// build
		{"Build Docker Image", "cd project && docker build -t innodeploy-app ."},

		// tag avec version
		{"Tag Image", "docker tag innodeploy-app " + imageName},

		// push
		{"Push Image", "docker push " + imageName},

		// 🔥 update YAML (GitOps)
		{"Update YAML", fmt.Sprintf(
			"sed -i 's|image:.*|image: %s|' k8s/deployment.yaml",
			imageName,
		)},

		// push git
		{"Git Commit", fmt.Sprintf(
			"git add . && git commit -m 'deploy %s'",
			version,
		)},
		{"Git Push", "git push"},
	}

	for _, step := range steps {
		runStep(step)
	}

	log.Println("✅ Pipeline finished 🚀 version:", version)
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