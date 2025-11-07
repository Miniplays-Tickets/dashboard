{#if modal}
    <StaffOverrideModal {guildId} on:close={() => modal = false} on:confirm={handleConfirm}/>
{/if}

<div class="parent">
    <div class="content">
        <div class="main-col">
            <Card footer footerRight>
                <span slot="title">
                    Support Zugriff
                </span>

                <div slot="body" class="body-wrapper">
                    Du kannst dem Support Team temporären Zugriff auf das Dashboard des Servers geben, sodass sie dir mit Problemen helfen können. 
                    Du kannst diesen Zugriff jederzeit auf dieser Webseite entfernen.
                </div>

                <div slot="footer" class="footer-wrapper">
                    {#if activeOverride}
                        <Button danger on:click={removeOverride}>
                            Zugriff Entfernen
                        </Button>
                    {/if}
                    <Button on:click={() => modal = true}>
                        Zugriff geben
                    </Button>
                </div>
            </Card>
        </div>
    </div>
</div>

<script>
    import Card from "../components/Card.svelte";
    import {notifyError, notifySuccess, withLoadingScreen} from '../js/util'
    import axios from "axios";
    import {API_URL} from "../js/constants";
    import {setDefaultHeaders} from '../includes/Auth.svelte'
    import Button from "../components/Button.svelte";
    import StaffOverrideModal from "../components/manage/StaffOverrideModal.svelte";

    export let currentRoute;
    let guildId = currentRoute.namedParams.id;

    let modal = false;
    let activeOverride = false;

    async function handleConfirm(e) {
        await createOverride(e.detail.timePeriod);
    }

    async function loadActiveOverride() {
        const res = await axios.get(`${API_URL}/api/${guildId}/staff-override`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        activeOverride = res.data.has_override;
    }

    async function createOverride(timePeriod) {
        let data = {
            time_period: timePeriod
        };

        const res = await axios.post(`${API_URL}/api/${guildId}/staff-override`, data);
        if (res.status !== 204) {
            notifyError(res.data.error);
            return;
        }

        modal = false;
        activeOverride = true;
        notifySuccess('Dem Support wurde der Zugriff gewährt');
    }

    async function removeOverride() {
        const res = await axios.delete(`${API_URL}/api/${guildId}/staff-override`);
        if (res.status !== 204) {
            notifyError(res.data.error);
            return;
        }

        activeOverride = false;
        notifySuccess('Dem Support wurde der Zugriff entfernt');
    }

    withLoadingScreen(async () => {
        setDefaultHeaders();
        await loadActiveOverride();
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
        justify-content: center;
        width: 100%;
        height: 100%;
        margin-top: 30px;
    }

    .main-col {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
    }

    .body-wrapper {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
    }

    .footer-wrapper {
        display: flex;
        flex-direction: row;
        gap: 10px;
        height: 100%;
    }

    @media only screen and (max-width: 950px) {
        .content {
            flex-direction: column-reverse;
        }

        .main-col {
            width: 100%;
            margin-top: 4%;
        }
    }
</style>

