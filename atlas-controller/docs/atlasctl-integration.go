// pkg/atlas/client.go - добавить поддержку AtlasApp CRD

func (c *Client) GetAtlasAppsFromCRD() ([]AtlasApp, error) {
    // Если установлен Atlas Controller, читаем из AtlasApp CRD
    atlasApps := []AtlasApp{}
    
    for _, namespace := range c.namespaces {
        // Попробовать получить AtlasApp ресурсы
        cmd := exec.Command("kubectl", "get", "atlasapp", 
            "-n", namespace, 
            "-o", "json",
            "--kubeconfig", c.kubeconfig)
        
        output, err := cmd.Output()
        if err != nil {
            // Если AtlasApp CRD не найден, fallback к старому методу
            continue
        }
        
        var atlasAppList AtlasAppList
        if err := json.Unmarshal(output, &atlasAppList); err != nil {
            continue
        }
        
        for _, item := range atlasAppList.Items {
            app := AtlasApp{
                Namespace:   item.Metadata.Namespace,
                App:         "atlas",
                Version:     item.Spec.Version,
                MigrationID: fmt.Sprintf("%d", item.Spec.MigrationId),
                Status:      item.Status.Phase,
                Replicas:    fmt.Sprintf("%d/%d", item.Status.ReadyReplicas, item.Status.TotalReplicas),
                LastUpdate:  item.Status.LastUpdate,
                Age:         calculateAge(item.Metadata.CreationTimestamp),
            }
            atlasApps = append(atlasApps, app)
        }
    }
    
    return atlasApps, nil
}

// Структуры для AtlasApp CRD
type AtlasAppList struct {
    Items []AtlasAppItem `json:"items"`
}

type AtlasAppItem struct {
    Metadata AtlasAppMetadata `json:"metadata"`
    Spec     AtlasAppSpec     `json:"spec"`
    Status   AtlasAppStatus   `json:"status"`
}

type AtlasAppMetadata struct {
    Namespace           string `json:"namespace"`
    CreationTimestamp   string `json:"creationTimestamp"`
}

type AtlasAppSpec struct {
    Version     string `json:"version"`
    MigrationId int    `json:"migrationId"`
}

type AtlasAppStatus struct {
    Phase         string `json:"phase"`
    ReadyReplicas int32  `json:"readyReplicas"`
    TotalReplicas int32  `json:"totalReplicas"`
    LastUpdate    string `json:"lastUpdate"`
}
