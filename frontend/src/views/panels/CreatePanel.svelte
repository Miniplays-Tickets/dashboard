<main>
    <a href="/manage/{guildId}/panels" class="link">
        <i class="fas fa-arrow-left"></i>
        Zurück zu Panels
    </a>
    <Card footer="{false}">
        <span slot="title">Panel erstellen</span>

        <div slot="body" class="body-wrapper">
            {#if !$loadingScreen}
                <PanelCreationForm {guildId} {channels} {roles} {emojis} {teams} {forms} {isPremium}
                                   bind:data={panelCreateData}/>
                <div class="submit-wrapper">
                    <Button icon="fas fa-paper-plane" fullWidth={true} on:click={createPanel}>Erstellen</Button>
                </div>
            {/if}
        </div>
    </Card>
</main>

<style>
    main {
        display: flex;
        flex-direction: column;
        width: 100%;
        row-gap: 1vh;
    }

    main > a {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 18px;
    }

    .body-wrapper {
        display: flex;
        flex-direction: column;
    }

    .submit-wrapper {
        margin: 1vh auto auto;
        width: 30%;
    }
</style>

<script>
    import {loadingScreen} from "../../js/stores";
    import Button from "../../components/Button.svelte";
    import Card from "../../components/Card.svelte";
    import PanelCreationForm from "../../components/manage/PanelCreationForm.svelte";
    import {setDefaultHeaders} from '../../includes/Auth.svelte'
    import {notifyError, setBlankStringsToNull, withLoadingScreen} from "../../js/util";
    import {onMount} from "svelte";
    import {loadChannels, loadEmojis, loadForms, loadPremium, loadRoles, loadTeams} from "../../js/common";
    import axios from "axios";
    import {API_URL} from "../../js/constants";
    import {Navigate, navigateTo} from "svelte-router-spa";

    setDefaultHeaders();

    export let currentRoute;
    let guildId = currentRoute.namedParams.id;

    let channels = [];
    let roles = [];
    let emojis = [];
    let teams = [];
    let forms = [];
    let isPremium = false;

    let panelCreateData;

    async function createPanel() {
        setBlankStringsToNull(panelCreateData);

        const res = await axios.post(`${API_URL}/api/${guildId}/panels`, panelCreateData);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        navigateTo(`/manage/${guildId}/panels?created=true`);
    }

    onMount(async () => {
        await withLoadingScreen(async () => {
            await Promise.all([
                loadChannels(guildId).then(r => channels = r).catch(e => notifyError(e)),
                loadRoles(guildId).then(r => roles = r).catch(e => notifyError(e)),
                loadEmojis(guildId).then(r => emojis = r).catch(e => notifyError(e)),
                loadTeams(guildId).then(r => teams = r).catch(e => notifyError(e)),
                loadForms(guildId).then(r => forms = r).catch(e => notifyError(e)),
                loadPremium(guildId, false).then(r => isPremium = r).catch(e => notifyError(e)),
            ])
        });
    });
</script>
