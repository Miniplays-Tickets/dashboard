<main>
    <Card fill={false} footer={false}>
        <span slot="title">Choose Premium Servers</span>
        <div slot="body" class="card-body">
            <div class="explanation">
                <span>
                    Dein Premium erlaubt dir {@html limit === 1 ? "<b>einen</b>" : `bis zu <b>${limit}</b>`} Premium Server.
                </span>
                <span>
                    Aktuell ausgewählt: <b>{selected.length} / {limit}</b> Server{limit > 1 ? "s" : ""}.
                </span>
            </div>
            <div class="servers">
                {#each getAdminGuilds(guilds) as guild}
                    <div class="server" class:active={selected.includes(guild.id)} class:pointer={selected.length < limit || selected.includes(guild.id)}
                         on:click={() => toggleSelected(guild.id)}>
                        <img src="{getIconUrl(guild.id, guild.icon)}" alt="Guild Icon" on:error={(e) => handleImgLoadError(e, guild.id)} />
                        <span class="name">{guild.name}</span>
                    </div>
                {/each}
            </div>
            <div class="submit-wrapper">
                <Button on:click={submitServers}>Speichern</Button>
            </div>
        </div>
    </Card>
</main>

<style>
    main {
        width: 100%;
        padding: 30px;
    }

    .card-body {
        display: flex;
        flex-direction: column;
        gap: 1em;
        padding-bottom: 1em;
        width: 100%;
    }

    .explanation {
        display: flex;
        flex-direction: column;
        gap: 1em;
    }

    .servers {
        display: flex;
        flex-wrap: wrap;
        gap: 1em;
        row-gap: 1em;
    }

    .server {
        display: flex;
        align-items: center;
        flex: 1 0 21%;
        gap: 1em;
        padding: 8px 10px;
        border-radius: 4px;
        user-select: none;
        background-color: #121212;
        border: 1px solid #121212;
        box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
    }

    .server.pointer {
        cursor: pointer;
    }

    .server.active {
        border-color: var(--primary);
        box-shadow: 0 0 10px var(--primary);
    }

    .server > img {
        width: 50px;
        height: 50px;
        border-radius: 50%;
    }

    .submit-wrapper {
        display: flex;
        justify-content: center;
    }
</style>

<script>
    import {withLoadingScreen, notifyError, notifySuccess, errorPage} from '../../js/util';
    import {setDefaultHeaders} from '../../includes/Auth.svelte'
    import Card from "../../components/Card.svelte";
    import {getIconUrl, getDefaultIcon} from "../../js/icons";
    import {API_URL} from "../../js/constants";
    import Button from "../../components/Button.svelte";
    import axios from "axios";

    let limit = 1;
    let selected = [];
    let guilds = [];

    function getAdminGuilds(guilds) {
        return guilds.filter(g => g.permission_level === 2);
    }

    let failed = [];
    function handleImgLoadError(e, guildId) {
        if (!failed.includes(guildId)) {
            failed = [...failed, guildId];
            e.target.src = getDefaultIcon(guildId);
        }
    }

    function toggleSelected(guildId) {
        if (selected.includes(guildId)) {
            selected = selected.filter(id => id !== guildId);
        } else {
            if (selected.length < limit) {
                selected = [...selected, guildId];
            }
        }
    }

    setDefaultHeaders();

    async function loadEntitlements() {
        const res = await axios.get(`${API_URL}/api/premium/@me/entitlements`)
        if (res.status !== 200) {
            notifyError(`Fehler beim Laden deiner Abonnements: ${res.data.error}`)
            return;
        }

        if (res.data.legacy_entitlement === null || res.data.legacy_entitlement.is_legacy) {
            errorPage('Diese Funktion ist nur für Benutzer mit einer bestimmten Premium Abonnement.');
            return;
        }

        limit = res.data.permitted_server_count;
        selected = res.data.selected_guilds;
    }

    async function loadGuilds() {
        const fromLocalStorage = window.localStorage.getItem('guilds');
        if (!fromLocalStorage) {
            notifyError('Fehler beim laden der Guilds aus dem lokalen Speicher.');
            return;
        }

        guilds = [...guilds, ...JSON.parse(fromLocalStorage)];
    }

    async function submitServers() {
        const res = await axios.put(`${API_URL}/api/premium/@me/active-guilds`, {
            selected_guilds: selected
        });

        if (res.status !== 204) {
            notifyError(`Fehler beim Speichern von Servern: ${res.data.error}`);
            return;
        }

        notifySuccess('Deine Premium Server Auswahl wurde gespeichert.')
    }

    withLoadingScreen(async () => {
        await Promise.all([
            loadEntitlements(),
            loadGuilds()
        ]);

        for (const id of selected) {
            if (!guilds.find(g => g.id === id)) {
                guilds = [{
                    id,
                    name: `Unbekannter Server ${id}`,
                    icon: "",
                    permission_level: 2
                }, ...guilds];
            }
        }
    });
</script>
