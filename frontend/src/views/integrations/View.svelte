<div class="parent">
  <div class="content">
    <div class="col-left">
      <Card footer footerRight>
        <span slot="title">Über {integration.name}</span>
        <div slot="body" class="body-wrapper">
          <span class="description">{integration.description}</span>

          <p style="padding-top: 5px">
            Wenn ein Benutzer ein Ticket öffnet, eine Anfrage, 
            welche die Benutzer ID des Ticket Öffners wird zu der folgenden URL gesendet, welche dem Integration Author gehört:
          </p>
          <input readonly value={integration.webhook_url} class="form-input"/>


          {#if privacy_policy_url === null}
            <p>Der Integration Author hat keine Privacy Policy zur Verfügung gestellt.</p>
          {:else}
            <p>
              Der Integration Author hat eine Privacy Policy zur Verfügung gestellt, abrufbar unter
              <a href="{privacy_policy_url}" class="link-blue">{privacy_policy_url}</a>
            </p>
          {/if}
        </div>
        <div slot="footer">
          {#if isActive}
            <Button on:click={removeIntegration} danger>Entfernen vom Server</Button>
          {:else}
            <Button on:click={() => navigateTo(`/manage/${guildId}/integrations/activate/${integrationId}`)}>
              Hinzufügen zum Server
            </Button>
          {/if}
        </div>
      </Card>
    </div>
    <div class="col-right">
      <Card footer={false} fill={false}>
        <span slot="title">Placeholders</span>
        <div slot="body">
          <p>Die folgenden Placeholder sind für den Benutzer in der Willkommensnachricht durch die <i>{integration.name}</i>
            Integration:</p>

          <div class="placeholders">
            {#if integration.placeholders}
              {#each integration.placeholders as placeholder}
                <Badge>%{placeholder.name}%</Badge>
              {/each}
            {/if}
          </div>
        </div>
      </Card>
    </div>
  </div>
</div>

<script>
    import {notifyError, notifySuccess, withLoadingScreen} from '../../js/util'
    import axios from "axios";
    import {API_URL} from "../../js/constants";
    import {setDefaultHeaders} from '../../includes/Auth.svelte'
    import Card from "../../components/Card.svelte";
    import Button from "../../components/Button.svelte";
    import Badge from "../../components/Badge.svelte";
    import {Navigate, navigateTo} from "svelte-router-spa";

    export let currentRoute;
    let guildId = currentRoute.namedParams.id;
    let integrationId = currentRoute.namedParams.integration;
    let freshlyCreated = currentRoute.queryParams.created === "true";

    let integration = {};
    let isActive = false;
    let privacy_policy_url = null;

    async function removeIntegration() {
        const res = await axios.delete(`${API_URL}/api/${guildId}/integrations/${integrationId}`);
        if (res.status !== 204) {
            notifyError(res.data.error);
            return;
        }

        navigateTo(`/manage/${guildId}/integrations?removed=true`);
    }

    async function loadIntegration() {
        let res = await axios.get(`${API_URL}/api/integrations/view/${integrationId}`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        integration = res.data;

        if (integration.privacy_policy_url !== null) {
            let tmp = new URL(integration.privacy_policy_url);
            if (tmp.protocol === "http:" || tmp.protocol === "https:") {
                privacy_policy_url = tmp;
            }
        }
    }

    async function loadIsActive() {
        let res = await axios.get(`${API_URL}/api/${guildId}/integrations/${integrationId}`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        isActive = res.data.active;
    }

    withLoadingScreen(async () => {
        setDefaultHeaders();

        await Promise.all([
            loadIntegration(),
            loadIsActive()
        ]);

        if (freshlyCreated) {
            notifySuccess("Deine Integration wurde erfolgreich erstellt!");
        }
    });
</script>

<style>
    .parent {
        display: flex;
        justify-content: center;
        width: 100%;
        height: 100%;
    }

    .content {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        width: 96%;
        height: 100%;
        margin-top: 30px;
        padding-bottom: 5vh;

        gap: 2%;
    }

    .col-left {
        width: 60%;
    }

    .col-right {
        width: 40%;
    }

    .placeholders {
        display: flex;
        flex-direction: row;
        flex-wrap: wrap;
        gap: 10px;
        margin-top: 10px;
    }

    .description {
        border-bottom: 1px solid #777;
        padding-bottom: 5px;
    }

    .body-wrapper {
        display: flex;
        flex-direction: column;
    }

    @media only screen and (max-width: 950px) {
        .content {
            flex-direction: column;
            row-gap: 2vh;
        }

        .col-left, .col-right {
            width: 100%;
        }
    }
</style>
