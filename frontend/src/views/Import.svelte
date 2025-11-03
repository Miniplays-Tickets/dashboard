<div class="content">
    <Card footer={false}>
        <span slot="title">Import von Einstellungen, Tickets & Transcripts</span>
        <div slot="body" class="body-wrapper">
            <div class="section">
                <h3 class="section-title">Import Items</h3>
                <form>
                    <div class="row">
                        <div class="col-4">
                            <label class="form-label" for="import_data">
                                Daten Export Datei (.zip)
                            </label>
                            <div class="col-1">
                                <input
                                    type="file"
                                    id="import_data"
                                    style="display: block; width: 90%;"
                                    accept=".zip"
                                />
                            </div>
                        </div>
                        <div class="col-4">
                            <label class="form-label" for="import_transcripts">
                                Transcripts Export Datei (.zip)
                            </label>
                            <div class="col-1">
                                <input
                                    type="file"
                                    id="import_transcripts"
                                    style="display: block; width: 90%;"
                                    accept=".zip"
                                />
                            </div>
                        </div>
                    </div>
                    <br />
                    <div class="row">
                        <div class="col-6">
                            <Button on:click={dispatchConfirm} icon={queryLoading ? "fa-solid fa-spinner fa-spin-pulse" : ""} disabled={queryLoading}>Speichern</Button>
                        </div>
                    </div>
                    <div class="row">
                        {#if queryLoading}
                        <div>
                            <br />
                            <br />
                            <p style="text-align: center;"><i class="fa-solid fa-spinner fa-spin-pulse"></i> Deine Daten werden gerade hochgeladen, bitte verlasse diese Seite nicht.</p>
                        </div>
                        {/if}
                    </div>
                </form>
            </div>

            {#if runs.length > 0}
            <div class="section">
                <h2 class="section-title">Ausführungen <span style="font-size: 12px; font-style: italic;">(Aktuallisiert alle 30 Sekunden)</span></h2>
                {#each ["DATA", "TRANSCRIPT"] as runType}
                    {#if runs.filter(run => run.run_type == runType).length > 0}
                        <h3>{runType.toLowerCase().replace(/\b\w/g, s => s.toUpperCase())} Logs</h3>
                        {#each runs.filter(run => run.run_type == runType).sort((a, b) => new Date(a.date) - new Date(b.date)) as run}
                        <Collapsible tooltip="Zeigt dir die logs für diese Ausführung">
                            <span slot="header" class="header">{run.run_type} Run #{run.run_id} - {new Date(run.date).toLocaleDateString('de-de', { weekday:"long", year:"numeric", month:"long", day:"numeric", hour: "2-digit", minute: "2-digit"})}</span>
                            <div slot="content" class="col-1">
                            <table class="nice">
                                <thead>
                                <tr>
                                    <th>Log Id</th>
                                    <th>Log Status</th>
                                    <th>Entrag Typ</th>
                                    <th>Nachricht</th>
                                    <th>Datum</th>
                                </tr>
                                </thead>
                                <tbody>
                                {#each run?.logs as log}
                                <tr>
                                    <td>{log.run_log_id}</td>
                                    <td>{log.log_type}</td>
                                    <td>{log.entity_type ?? "N/A"}</td>
                                    <td>{log.message ?? "N/A"}</td>
                                    <td>{new Date(log.date).toLocaleDateString('de-de', { weekday:"long", year:"numeric", month:"long", day:"numeric", hour: "2-digit", minute: "2-digit", second: "2-digit"})}</td>
                                </tr>
                                {/each}
                                </tbody>
                            </table>
                            </div>
                        </Collapsible>
                        {/each}
                    {/if}
                {/each}
            </div>
            {/if}

            {#if dataReturned}
            <div class="section">
                <h2 class="section-title">Import Dateien Hochgeladen</h2>
                <div class="row">
                    <p style="text-align: center;">Deine Daten & Transcripts wurden der Wartenschlange hinzugefügt, es wird etwas dauern, bis diese Importiert wurden.</p>
                </div>
            </div>
            {/if}
        </div>
    </Card>
</div>
<script>
    import { createEventDispatcher } from "svelte";
    import { fade } from "svelte/transition";
    import Card from "../components/Card.svelte";
    import Button from "../components/Button.svelte";

    import Textarea from "../components/form/Textarea.svelte";

    import { setDefaultHeaders } from "../includes/Auth.svelte";
    import { notify, notifyError, notifySuccess } from "../js/util";
    import axios from "axios";
    import { API_URL } from "../js/constants";
    import Collapsible from "../components/Collapsible.svelte";
    setDefaultHeaders();

    export let currentRoute;
    let guildId = currentRoute.namedParams.id

    let dataReturned = false;

    let queryLoading = false;

    let runs = [];

    let queuePositions = {
        data: 0,
        transcripts: 0,
    };

    const dispatch = createEventDispatcher();

    function dispatchClose() {
        dispatch("close", {});
    }

    axios.get(`${API_URL}/api/${guildId}/import/queue`).then((res) => {
        if (res.status !== 200) {
            notifyError(`Failed to get import queue: ${res.data.error}`);
            return;
        }

        queuePositions = res.data;
    });

    function getRuns() {
        axios.get(`${API_URL}/api/${guildId}/import/runs`).then((res) => {
            if (res.status !== 200) {
                notifyError(`Fehler beim abrufen von Ausführungen: ${res.data.error}`);
                return;
            }

            runs = res.data;
        }); 
    }

    getRuns();

    setInterval(() => {
        getRuns();
    }, 30 * 1000);


    async function dispatchConfirm() {
        let dataFileInput = document.getElementById("import_data");
        let transcriptFileInput = document.getElementById("import_transcripts");

        if (
            dataFileInput.files.length === 0 &&
            transcriptFileInput.files.length === 0
        ) {
            notifyError(
                "Bitte wähle eine Datei zum importieren aus.",
            );
            return;
        }

        const frmData = new FormData();
        if (dataFileInput.files.length > 0) {
            frmData.append("data_file", dataFileInput.files[0]);
        }

        queryLoading = true;
        setTimeout(() => {
            if (queryLoading) {
                notify(
                    "Hochladen...",
                    "Deine Dateien werden noch Hochgeladen, bitte warte kurz.",
                );
            }
        }, 60 * 1000);

        if (transcriptFileInput.files.length > 0) {
            const presignTranscriptRes = await axios.get(`${API_URL}/api/${guildId}/import/presign?file_size=${transcriptFileInput.files[0].size}&file_type=transcripts&file_content_type=${transcriptFileInput.files[0].type}`);
            if (presignTranscriptRes.status !== 200) {
                notifyError(`Fehler beim Hochladen der Transscripts: ${presignTranscriptRes.data.error}`);
                queryLoading = false;
                return;
            }
            
            await fetch(presignTranscriptRes.data.url, {
                method: "PUT",
                body: transcriptFileInput.files[0],
                headers: {
                    "Content-Type": transcriptFileInput.files[0].type,
                },
            }).then((res) => {
                if (res.status !== 200) {
                    notifyError(`Fehler beim Hochladen der Transscripts: ${res.data.error}`);
                    queryLoading = false;
                    return;
                }

                dataReturned = true;
                notifySuccess(`Transcripts uploaded successfully - There are currently ${queuePositions.transcripts} import${queuePositions.transcripts == 1 ? "" : "s"} ahead of you in the transcript queue. These can take a while to process, please check back later`);
            }).catch((e) => {
                notifyError(`Failed to upload transcripts: ${e}`);
                queryLoading = false;
            });
        }

        if (dataFileInput.files.length > 0) {
            const presignDataRes = await axios.get(`${API_URL}/api/${guildId}/import/presign?file_size=${dataFileInput.files[0].size}&file_type=data&file_content_type=${dataFileInput.files[0].type}`);
            if (presignDataRes.status !== 200) {
                notifyError(`Fehler beim Hochladen der Serverdaten: ${presignDataRes.data.error}`);
                queryLoading = false;
                return;
            }
            
            await fetch(presignDataRes.data.url, {
                method: "PUT",
                body: dataFileInput.files[0],
                headers: {
                    "Content-Type": dataFileInput.files[0].type,
                },
            }).then((res) => {
                if (res.status !== 200) {
                    notifyError(`Fehler beim Hochladen der Serverdaten: ${res.data.error}`);
                    queryLoading = false;
                    return;
                }

                dataReturned = true;
                notifySuccess(`Data uploaded successfully - There are currently ${queuePositions.data} import${queuePositions.data == 1 ? "" : "s"} ahead of you in the data queue. These can take a while to process, please check back later`);
            }).catch((e) => {
                notifyError(`Failed to upload data: ${e}`);
                queryLoading = false;
            });
        }

        queryLoading = false;

        dispatchClose();
    }

    function handleKeydown(e) {
        if (e.key === "Escape") {
            dispatchClose();
        }
    }
    
</script>
<style>
    .content {
        display: flex;
        width: 100%;
        height: 100%;
    }

    .body-wrapper {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
        padding: 1%;
    }

    .section {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
    }

    .section:not(:first-child) {
        margin-top: 2%;
    }

    .section-title {
        font-size: 36px;
        font-weight: bolder !important;
    }

    h3 {
        font-size: 28px;
        margin-bottom: 4px;
    }

    .row {
        display: flex;
        flex-direction: row;
        width: 100%;
        height: 100%;
    }

    ul {
        margin: 0;
        padding: 0;
    }

    li {
        list-style-type: none;
    }

    .manage {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        width: 100%;
        height: 100%;
        margin-top: 12px;
    }

    .col-4 {
        width: 100% !important;
    }

    .col {
        display: flex;
        flex-direction: column;
        align-items: center;
        width: 49%;
        height: 100%;
    }

    table.nice > tbody > tr:first-child {
        border-top: 1px solid #dee2e6;
    }

    .user-select {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        width: 100%;
        height: 100%;
        margin-bottom: 1%;
    }

    @media only screen and (max-width: 950px) {
        .manage {
            flex-direction: column;
        }

        .col {
            width: 100%;
        }
    }
</style>
