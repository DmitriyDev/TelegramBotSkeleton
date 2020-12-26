--- 
#Skeleton for telegram bot

##Deployment

#### Step 1. Create and fill .env (from .env.dist)
#### Step 2. Run docker-compose
```docker-compose up -d```

---

## Add new endpoint for command

#### Step 1. Add handler function to bot_processor

```
func (bp *BotProcessor) NewEndpointHandler(m *tb.Message) {
	bp.Send(m.Sender, "Hello new message")
}

```

#### Step 2. Update ```config.yaml``` (add new config for endpoints)

```

  - name: newCommand
    command: /newCommand
    description: New command description
    handler: NewEndpointHandler
    admin: false
    visible: true

```

```
name        Command name
command     Command ( "/" required)
description New command description
handler     Name of method created in step 1
admin       (true/false) If true, this command would be not allowed for non-admin user
visible     (true/false) Allowed to execute, but not displayed on /start/help/hello command responses
```
