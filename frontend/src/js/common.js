import axios from "axios";
import {API_URL} from "./constants";

export async function loadPremium(guildId, includeVoting = false) {
    const res = await axios.get(`${API_URL}/api/${guildId}/premium?include_voting=${includeVoting}`);
    if (res.status !== 200) {
        throw new Error(`Fehler beim Laden vom Premiumstatus: ${res.data.error}`);
    }

    return res.data.premium;
}

export async function loadChannels(guildId) {
    const res = await axios.get(`${API_URL}/api/${guildId}/channels`);
    if (res.status !== 200) {
        throw new Error(`Fehler beim laden von Kanalen: ${res.data.error}`);
    }

    return res.data;
}

export async function loadPanels(guildId) {
    const res = await axios.get(`${API_URL}/api/${guildId}/panels`);
    if (res.status !== 200) {
        throw new Error(`Fehler beim laden von den Panels: ${res.data.error}`);
    }

    // convert button_style and form_id to string
    return res.data.map((p) => Object.assign({}, p, {
        button_style: p.button_style.toString(),
        form_id: p.form_id === null ? "null" : p.form_id
    }));
}

export async function loadMultiPanels(guildId) {
    const res = await axios.get(`${API_URL}/api/${guildId}/multipanels`);
    if (res.status !== 200) {
        throw new Error(`Fehler beim laden von den Multi-Panels: ${res.data.error}`);
    }

    return res.data.data;
}

export async function loadTeams(guildId) {
    const res = await axios.get(`${API_URL}/api/${guildId}/team`);
    if (res.status !== 200) {
        throw new Error(`Fehler beim Laden von Teams: ${res.data.error}`);
    }

    return res.data;
}

export async function loadRoles(guildId) {
    const res = await axios.get(`${API_URL}/api/${guildId}/roles`);
    if (res.status !== 200) {
        throw new Error(`Fehler beim laden von Rollen: ${res.data.error}`);
    }

    return res.data.roles;
}

export async function loadEmojis(guildId) {
    const res = await axios.get(`${API_URL}/api/${guildId}/emojis`);
    if (res.status !== 200) {
        throw new Error(`Fehler beim laden vom Emojis: ${res.data.error}`);
    }

    return res.data;
}

export async function loadForms(guildId) {
    const res = await axios.get(`${API_URL}/api/${guildId}/forms`);
    if (res.status !== 200) {
        throw new Error(`Fehler beim Laden von Formularen: ${res.data.error}`);
    }

    return res.data || [];
}
