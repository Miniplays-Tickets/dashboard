{#if panelToDelete !== null}
    <ConfirmationModal icon="fas fa-trash-can" isDangerous on:cancel={() => panelToDelete = null}
                       on:confirm={() => deletePanel(panelToDelete.panel_id)}>
        <span slot="body">Bist du dir Sicher dass du das Panel {panelToDelete.title} löschen möchtest?</span>
        <span slot="confirm">Löschen</span>
    </ConfirmationModal>
{/if}

{#if multiPanelToDelete !== null}
    <ConfirmationModal icon="fas fa-trash-can" isDangerous on:cancel={() => multiPanelToDelete = null}
                       on:confirm={() => deleteMultiPanel(multiPanelToDelete.id)}>
        <span slot="body">Bist du dir Sicher dass du dass Multi-Panel {multiPanelToDelete.embed?.title || "Öffne ein Ticket!!"} löschen möchtest?</span>
        <span slot="confirm">Löschen</span>
    </ConfirmationModal>
{/if}

<div class="wrapper">
    <div class="col">
        <div class="row">
            <Card footer="{false}">
                <span slot="title">Ticket Panel</span>
                <div slot="body" class="card-body panels">
                    <div class="controls">
                        <p>Deine Panel Anzahl: <b>{panels.length} / {isPremium ? '∞' : '3'}</b></p>
                        <Navigate to="/manage/{guildId}/panels/create" styles="link">
                            <Button icon="fas fa-plus">Neues Panel</Button>
                        </Navigate>
                    </div>

                    <table style="margin-top: 10px">
                        <thead>
                        <tr>
                            <th>Kanal</th>
                            <th class="max">Panel Titel</th>
                            <th></th>
                            <th></th>
                            <th></th>
                        </tr>
                        </thead>
                        <tbody>
                        {#each panels as panel}
                            <tr>
                                <td>#{channels.find((c) => c.id === panel.channel_id)?.name ?? 'Unbekannter Kanal'}</td>
                                <td class="max">{panel.title}</td>
                                <td>
                                    <Button disabled={panel.force_disabled}
                                            on:click={() => resendPanel(panel.panel_id)}>Neu Senden
                                    </Button>
                                </td>
                                <td>
                                    <Navigate to="/manage/{guildId}/panels/edit/{panel.panel_id}" styles="link">
                                        <Button disabled={panel.force_disabled}>Bearbeiten</Button>
                                    </Navigate>
                                </td>
                                <td>
                                    <Button danger on:click={() => panelToDelete = panel}>Löschen</Button>
                                </td>
                            </tr>
                        {/each}
                        </tbody>
                    </table>
                </div>
            </Card>
        </div>
    </div>
    <div class="col">
        <div class="row">
            <Card footer="{false}">
                <span slot="title">Multi-Panels</span>
                <div slot="body" class="card-body">
                    <div class="controls">
                        <Navigate to="/manage/{guildId}/panels/create-multi" styles="link">
                            <Button icon="fas fa-plus">Neues Multi-Panel</Button>
                        </Navigate>
                    </div>

                    <table style="margin-top: 10px">
                        <thead>
                        <tr>
                            <th class="max">Panel Titel</th>
                            <th></th>
                            <th></th>
                            <th></th>
                        </tr>
                        </thead>
                        <tbody>
                        {#each multiPanels as panel}
                            <tr>
                                <td class="max">{panel.embed?.title || 'Öffne ein Ticket!'}</td>
                                <td>
                                    <Button on:click={() => resendMultiPanel(panel.id)}>Neu Senden</Button>
                                </td>
                                <td>
                                    <Navigate to="/manage/{guildId}/panels/edit-multi/{panel.id}" styles="link">
                                        <Button>Bearbeiten</Button>
                                    </Navigate>
                                </td>
                                <td>
                                    <Button danger on:click={() => multiPanelToDelete = panel}>Löschen</Button>
                                </td>
                            </tr>
                        {/each}
                        </tbody>
                    </table>
                </div>
            </Card>
        </div>
        <div class="row">

        </div>
    </div>
</div>

<script>
    import Card from "../../components/Card.svelte";
    import {checkForParamAndRewrite, notifyError, notifySuccess, withLoadingScreen} from "../../js/util";
    import axios from "axios";
    import {API_URL} from "../../js/constants";
    import {setDefaultHeaders} from '../../includes/Auth.svelte'
    import Button from "../../components/Button.svelte";
    import ConfirmationModal from "../../components/ConfirmationModal.svelte";
    import {Navigate} from "svelte-router-spa";
    import {loadChannels, loadMultiPanels, loadPanels, loadPremium} from "../../js/common";

    export let currentRoute;

    setDefaultHeaders()

    let guildId = currentRoute.namedParams.id;

    let channels = [];
    let panels = [];
    let multiPanels = [];
    let isPremium = false;

    let panelToDelete = null;
    let multiPanelToDelete = null;

    async function resendPanel(panelId) {
        const res = await axios.post(`${API_URL}/api/${guildId}/panels/${panelId}`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        notifySuccess("Panel erfolgreich neu gesendet");
    }

    async function deletePanel(panelId) {
        const res = await axios.delete(`${API_URL}/api/${guildId}/panels/${panelId}`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        panels = panels.filter((p) => p.panel_id !== panelId);
        panelToDelete = null;
    }

    async function resendMultiPanel(id) {
        const res = await axios.post(`${API_URL}/api/${guildId}/multipanels/${id}`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        notifySuccess("Multipanel erfolgreich neu gesendet")
    }

    async function deleteMultiPanel(id) {
        const res = await axios.delete(`${API_URL}/api/${guildId}/multipanels/${id}`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        multiPanels = multiPanels.filter((p) => p.id !== id);
        multiPanelToDelete = null;
    }

    withLoadingScreen(async () => {
        await Promise.all([
            loadChannels(guildId).then(r => channels = r).catch(e => notifyError(e)),
            loadPremium(guildId, false).then(r => isPremium = r).catch(e => notifyError(e)),
            loadPanels(guildId).then(r => panels = r).catch(e => notifyError(e)),
            loadMultiPanels(guildId).then(r => multiPanels = r).catch(e => notifyError(e))
        ])

        if (checkForParamAndRewrite("created")) {
            notifySuccess("Panel erfolgreich erstellt");
        }

        if (checkForParamAndRewrite("edited")) {
            notifySuccess("Panel erfolgreich bearbeitet");
        }

        if (checkForParamAndRewrite("notfound")) {
            notifyError("Panel nicht gefunden");
        }
    });
</script>

<style>
    .wrapper {
        display: flex;
        flex-direction: row;
        height: 100%;
        width: 100%;
        gap: 2%;
    }

    .col {
        display: flex;
        flex-direction: column;
        align-items: center;
        width: 50%;
    }

    .row {
        display: flex;
        width: 100%;
        margin-bottom: 2%;
    }

    .card-body {
        width: 100%;
    }

    .card-body.panels {
        display: flex;
        flex-direction: column;
        row-gap: 4%;
    }

    .card-body > .controls {
        display: flex;
        justify-content: right;
        align-items: center;
        gap: 2%;
    }

    .card-body.panels > .controls {
        justify-content: space-between;
    }

    @media only screen and (max-width: 1400px) {
        .wrapper {
            flex-direction: column;
        }

        .col {
            width: 100%;
        }
    }

    @media only screen and (max-width: 576px) {
        .row {
            width: 100%;
        }
    }

    table {
        width: 100%;
        border-collapse: collapse;
    }

    th {
        text-align: left;
        font-weight: normal;
        border-bottom: 1px solid #dee2e6;
        padding-left: 10px;
    }

    tr {
        border-bottom: 1px solid #dee2e6;
    }

    tr:last-child {
        border-bottom: none;
    }

    td {
        padding: 10px;
    }

    th {
        padding: 0 10px;
    }

    th:not(.max), td:not(.max) {
        width: 0;
        white-space: nowrap;
    }
</style>
