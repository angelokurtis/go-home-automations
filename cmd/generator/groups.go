package main

type Groups []Group

type Group struct {
	Switches []string
}

var groups Groups = Groups{
	Group{Switches: []string{
		"switch.interruptor_6x_da_copa_l3",
		"switch.arvore_de_natal",
	}},
	Group{Switches: []string{
		"switch.interruptor_6x_da_copa_l1",
		"switch.interruptor_6x_da_copa_l2",
		"switch.interruptor_3x_da_entrada_left",
	}},
	Group{Switches: []string{
		"switch.interruptor_3x_da_entrada_center",
		"switch.interruptor_6x_da_entrada_l1",
	}},
	Group{Switches: []string{
		"switch.interruptor_3x_da_entrada_right",
		"switch.interruptor_6x_da_entrada_l2",
	}},
	Group{Switches: []string{
		"switch.interruptor_2x_da_area_gourmet_center",
		"switch.interruptor_6x_da_area_gourmet_l5",
		"switch.interruptor_6x_da_copa_l4",
	}},
	Group{Switches: []string{
		"switch.interruptor_6x_da_area_gourmet_l1",
		"switch.interruptor_6x_da_area_gourmet_l6",
	}},
	Group{Switches: []string{
		"switch.interruptor_2x_da_area_gourmet_left",
		"switch.interruptor_6x_da_copa_l5",
	}},
	Group{Switches: []string{
		"switch.interruptor_6x_da_entrada_l3",
		"switch.interruptor_6x_da_entrada_l4",
	}},
	Group{Switches: []string{
		"switch.interruptor_1x_do_mezanino",
		"switch.interruptor_6x_da_copa_l6",
	}},
	Group{Switches: []string{
		"switch.interruptor_da_cozinha_l1",
		"switch.interruptor_da_cozinha2_l1",
	}},
	Group{Switches: []string{
		"switch.interruptor_da_lavanderia_l1",
		"switch.interruptor_da_cozinha_l2",
		"switch.interruptor_da_cozinha2_l2",
	}},
	Group{Switches: []string{
		"switch.interruptor_6x_da_copa2_l1",
		"switch.interruptor_6x_da_copa2_l2",
	}},
}
