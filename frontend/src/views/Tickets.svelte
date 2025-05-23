<main>
    <Card footer={false}>
        <span slot="title">
            <i class="fas fa-filter"></i>
            Filter
        </span>
        <div slot="body" class="filter-wrapper">
            <div>
                <label class="form-label">Zeilen zum Anzeigen</label>
                <ColumnSelector
                    options={["ID", "Panel", "Benutzer", "Öffnungszeit", "Beansprucht von", "Letzte Nachricht", "Wartet auf Antwort"]}
                    bind:selected={selectedColumns}
                />
            </div>

            <Dropdown col2 label="Sortiere Tickets in..." bind:value={sortMethod}>
                <option value="id_asc">Ticket ID (Steigend) / Älteste zuerst</option>
                <option value="id_desc">Ticket ID (Absteigend) / Neuste zuerst</option>
                <option value="unclaimed">Nicht Beansprucht & Wartet auf Antwort</option>
            </Dropdown>

            <Checkbox label="Zeige nur unbeanspruchte Tickets & Tickets von mir Beansprucht" bind:value={onlyShowMyTickets} />
        </div>
    </Card>

    <Card footer={false}>
        <span slot="title">Offene Tickets</span>
        <div slot="body" class="body-wrapper">
            <table class="nice">
                <thead>
                <tr>
                    <th class:visible={selectedColumns.includes('ID')}>ID</th>
                    <th class:visible={selectedColumns.includes('Panel')}>Panel</th>
                    <th class:visible={selectedColumns.includes('Benutzer')}>Benuter</th>
                    <th class:visible={selectedColumns.includes('Öffnungszeit')}>Öffnungszeit</th>
                    <th class:visible={selectedColumns.includes('Beansprucht von')}>Beansprucht von</th>
                    <th class:visible={selectedColumns.includes('Letzte Nachricht')}>Letzte Nachricht</th>
                    <th class:visible={selectedColumns.includes('Wartet auf Antwort')}>Wartet auf Antwort</th>
                    <th class="visible">Anzeigen</th>
                </tr>
                </thead>
                <tbody>
                {#each filtered as ticket}
                    {@const user = data.resolved_users[ticket.user_id]}
                    {@const claimer = ticket.claimed_by ? data.resolved_users[ticket.claimed_by] : null}
                    {@const panel_title = data.panel_titles[ticket.panel_id?.toString()]}

                    <tr>
                        <td class:visible={selectedColumns.includes('ID')}>{ticket.id}</td>
                        <td class:visible={selectedColumns.includes('Panel')}>
                            {panel_title || 'Unbekanntes Panel'}
                        </td>

                        <td class:visible={selectedColumns.includes('Benutzer')}>
                            {#if user}
                                {user.global_name || user.username}
                            {:else}
                                Unbekannt
                            {/if}
                        </td>

                        <td class:visible={selectedColumns.includes('Öffnungszeit')}>
                            {getRelativeTime(new Date(ticket.opened_at))}
                        </td>

                        <td class:visible={selectedColumns.includes('Beansprucht von')}>
                            {#if ticket.claimed_by === null}
                                <b>Nicht Beansprucht</b>
                            {:else if claimer}
                                {claimer.global_name || claimer.username}
                            {:else}
                                Unknown
                            {/if}
                        </td>

                        <td class:visible={selectedColumns.includes('Letzte Nachricht')}>
                            {#if ticket.last_response_time}
                                {getRelativeTime(new Date(ticket.last_response_time))}
                            {:else}
                                Nie
                            {/if}
                        </td>

                        <td class:visible={selectedColumns.includes('Wartet auf Antwort')}>
                            {#if ticket.last_response_is_staff}
                                Nein
                            {:else}
                                <b>Ja</b>
                            {/if}
                        </td>

                        <td class="visible">
                            <Navigate to="/manage/{guildId}/tickets/view/{ticket.id}" styles="link">
                                <Button type="button">Anzeigen</Button>
                            </Navigate>
                        </td>
                    </tr>
                {/each}
                </tbody>
            </table>
        </div>
    </Card>
</main>

<script>
    import Card from "../components/Card.svelte";
    import {getRelativeTime, notifyError, withLoadingScreen} from '../js/util'
    import axios from "axios";
    import {API_URL} from "../js/constants";
    import {setDefaultHeaders} from '../includes/Auth.svelte'
    import Button from "../components/Button.svelte";
    import {Navigate} from 'svelte-router-spa';
    import ColumnSelector from "../components/ColumnSelector.svelte";
    import Dropdown from "../components/form/Dropdown.svelte";
    import Checkbox from "../components/form/Checkbox.svelte";

    export let currentRoute;
    let guildId = currentRoute.namedParams.id;

    let selectedColumns = ['ID', 'Panel', 'User', 'Claimed By', 'Last Message Time', 'Awaiting Response'];
    let sortMethod = "unclaimed";
    let onlyShowMyTickets = false;

    let data = {
        tickets: [],
        panel_titles: {},
        resolved_users: {}
    };

    let filtered = [];

    function filterTickets() {
        filtered = data.tickets.filter(ticket => {
            if (onlyShowMyTickets) {
                return ticket.claimed_by === null || ticket.claimed_by === data.self_id;
            }

            return true;
        });

        // Apply sort
        if (sortMethod === "id_asc") {
            filtered.sort((a, b) => a.id - b.id);
        } else if (sortMethod === "id_desc") {
            filtered.sort((a, b) => b.id - a.id);
        } else if (sortMethod === 'unclaimed') {
            filtered.sort((a, b) => {
                // Place unclaimed tickets at the top. The priority of fields used for sorting is:
                // 1. Unclaimed tickets, or tickets claimed by the current user
                // 2. Awaiting Response
                // 3. Last Response Time

                // Unclaimed tickets at the top
                if (a.claimed_by === null && b.claimed_by !== null) {
                    return -1;
                }
                if (a.claimed_by !== null && b.claimed_by === null) {
                    return 1;
                }

                if (a.claimed_by === data.self_id && b.claimed_by !== data.self_id) {
                    return -1;
                }
                if (a.claimed_by !== data.self_id && b.claimed_by === data.self_id) {
                    return 1;
                }

                // Among claimed tickets, those awaiting response at the top
                if (!a.last_response_is_staff && b.last_response_is_staff) {
                    return -1;
                }
                if (a.last_response_is_staff && !b.last_response_is_staff) {
                    return 1;
                }

                // Among tickets not awaiting response, sort by last response time
                const aLastResponseTime = new Date(a.last_response_time || 0);
                const bLastResponseTime = new Date(b.last_response_time || 0);

                return aLastResponseTime - bLastResponseTime;
            });
        }
    }

    async function loadTickets() {
        const res = await axios.get(`${API_URL}/api/${guildId}/tickets`);
        if (res.status !== 200) {
            notifyError(res.data.error);
            return;
        }

        data = res.data;

        if (!data.tickets) {
            data.tickets = [];
        }

        data.tickets = data.tickets.map(ticket => {
            if (ticket.claimed_by === "null") {
                ticket.claimed_by = null;
            }

            return ticket;
        })

        filterTickets();
    }

    const columnStorageKey = 'ticket_list:selected_columns';
    const sortOrderKey = 'ticket_list:sort_order';
    const onlyMyTicketsKey = 'ticket_list:only_my_tickets';

    $: selectedColumns, updateFilters();
    $: sortMethod, updateFilters();
    $: onlyShowMyTickets, updateFilters();

    function updateFilters() {
        window.localStorage.setItem(columnStorageKey, JSON.stringify(selectedColumns));
        window.localStorage.setItem(sortOrderKey, sortMethod);
        window.localStorage.setItem(onlyMyTicketsKey, JSON.stringify(onlyShowMyTickets));

        filterTickets();
    }

    function loadFilterSettings() {
        const columns = window.localStorage.getItem(columnStorageKey);
        if (columns) {
            selectedColumns = JSON.parse(columns);
        }

        const sortOrder = window.localStorage.getItem(sortOrderKey);
        if (sortOrder) {
            sortMethod = sortOrder;
        }

        const onlyMyTickets = window.localStorage.getItem(onlyMyTicketsKey);
        if (onlyMyTickets) {
            onlyShowMyTickets = JSON.parse(onlyMyTickets);
        }
    }

    withLoadingScreen(async () => {
        loadFilterSettings();

        setDefaultHeaders();
        await loadTickets();
    });
</script>

<style>
    main {
        display: flex;
        flex-direction: column;
        gap: 30px;
        width: 100%;
        height: 100%;
    }

    .body-wrapper {
        display: flex;
        flex-direction: column;
        width: 100%;
        height: 100%;
    }

    .filter-wrapper {
        display: flex;
        flex-direction: row;
        gap: 1rem;
        width: 100%;
        height: 100%;
    }

    th, td {
        display: none;
    }

    th.visible, td.visible {
        display: table-cell;
    }

    @media only screen and (max-width: 1400px) {
        .filter-wrapper {
            flex-direction: column;
            gap: 8px;
        }
    }
</style>
