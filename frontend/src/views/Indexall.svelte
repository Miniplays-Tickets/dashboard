<div class="content">
  <div class="card-wrapper">
    <Card footer={false} fill={false}>
      <span slot="title">
        Admin Serverübersicht
      </span>

      <div slot="body" style="width: 100%">
        <div id="guild-container">
          {#each guildsall as guild}
            <Guild guild={guild}/>
          {/each}
        </div>

        <div class="flex-container" id="refresh-container">
          <Button icon="fas fa-sync" on:click={refreshGuilds}>
            Liste Aktuallisieren
          </Button>
        </div>
      </div>
    </Card>
  </div>
</div>

<script>
    import axios from 'axios';
    import {fade} from 'svelte/transition';
    import {notifyError, withLoadingScreen} from '../js/util'
    import {setDefaultHeaders} from '../includes/Auth.svelte'
    import {API_URL} from "../js/constants.js";
    import Guild from '../components/Guild.svelte'
    import Card from '../components/Card.svelte'
    import InviteBadge from '../components/InviteBadge.svelte'
    import Button from '../components/Button.svelte'
    import {loadingScreen, permissionLevelCache} from "../js/stores";

    setDefaultHeaders();

    let guildsall = window.localStorage.getItem('guildsall') ? JSON.parse(window.localStorage.getItem('guildsall')) : [];

    async function refreshGuilds() {
        await withLoadingScreen(async () => {
            const res = await axios.post(`${API_URL}/user/guilds/reloadall`);
            if (res.status !== 200) {
                notifyError(res.data.error);
                return;
            }

            if (!res.data.success && res.data['reauthenticate_required'] === true) {
                window.location.href = "/login";
                return;
            }

            guildsall = res.data.guilds;
            window.localStorage.setItem('guildsall', JSON.stringify(guildsall));
        });
    }

    loadingScreen.set(false);
</script>

<style>
    .content {
        display: flex;
        height: 100%;
        width: 100%;
        justify-content: center;
    }

    .card-wrapper {
        display: block;
        width: 75%;
        margin-top: 5%;
    }

    #guild-container {
        display: flex;
        flex-direction: row;
        flex-wrap: wrap;
        justify-content: space-evenly;
    }

    #refresh-container {
        display: flex;
        justify-content: center;

        margin: 10px 0;
        color: white;
    }

    @media (max-width: 576px) {
        .card-wrapper {
            width: 100%;
        }
    }
</style>
