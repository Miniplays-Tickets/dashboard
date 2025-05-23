<Card footer="{false}" fill="{false}">
  <span slot="title">
    Beanspruchungseinstellungen
  </span>

  <div slot="body" class="body-wrapper">
    <form class="settings-form" on:submit|preventDefault={updateSettings}>
      <div class="row">
        <Checkbox label="SUPPORT MITGLIEDER KÖNNEN BEANSPRUCHTE TICKETS ANSEHEN" col2=true bind:value={data.support_can_view} on:change={validateView} />
        <Checkbox label="SUPPORT MITGLIEDER KÖNNEN IN BEANSPRUCHTE TICKETS SCHREIBEN" col2=true bind:value={data.support_can_type} on:change={validateType} />
      </div>
      <div class="row">
        <div class="col-1">
          <Button icon="fas fa-paper-plane" fullWidth=true>Speichern</Button>
        </div>
      </div>
    </form>
  </div>
</Card>

<script>
    export let guildId;

    import Card from "../Card.svelte";
    import Checkbox from "../form/Checkbox.svelte";

    import {setDefaultHeaders} from '../../includes/Auth.svelte'
    import axios from "axios";
    import {notifyError, notifySuccess, withLoadingScreen} from "../../js/util";
    import {API_URL} from "../../js/constants";
    import Button from "../Button.svelte";

    setDefaultHeaders();

    let data = {
        support_can_view: true,
        support_can_type: true,
    };

    function validateView() {
        if (!data.support_can_view && data.support_can_type) {
            data.support_can_type = false;
        }
    }

    function validateType() {
        if (!data.support_can_view && data.support_can_type) {
            data.support_can_view = true;
        }
    }

    async function updateSettings() {
        const res = await axios.post(`${API_URL}/api/${guildId}/claimsettings`, data);
        if (res.status === 200 && res.data.success) {
            notifySuccess("Deine Einstellungen wurden gespeichert.");
        } else {
            notifyError(res.data.error);
        }
    }

    async function loadData() {
        const res = await axios.get(`${API_URL}/api/${guildId}/claimsettings`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        data = res.data;
    }

    withLoadingScreen(async () => {
        await loadData();
    });
</script>

<style>
    .body-wrapper {
        display: flex;
        width: 100%;
        height: 100%;
    }

    .row {
        display: flex;
        justify-content: space-between;
        width: 100%;
        height: 100%;
    }

    .settings-form {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
    }

    @media only screen and (max-width: 900px) {
        .row {
            flex-direction: column;
            justify-content: center;
        }
    }
</style>